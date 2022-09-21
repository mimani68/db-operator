package sql

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	corev1 "k8s.io/api/core/v1"
)

var (
	errorConnectionStates = []string{
		"connecting to master",
		"reconnecting after a failed binlog dump request",
		"reconnecting after a failed master event read",
		"waiting to reconnect after a failed binlog dump request",
		"waiting to reconnect after a failed master event read",
	}
)

type SQLRunner struct {
	db *sql.DB
}

func NewSQLRunner(user, password, host string, port int) (*SQLRunner, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/?timeout=5s&interpolateParams=true&multiStatements=true",
		user, password, host, port,
	)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &SQLRunner{db}, nil
}

func (s *SQLRunner) CheckSlaveStatusWithRetry(retry uint32) (isLagged, isReplicating corev1.ConditionStatus, err error) {
	for {
		if retry == 0 {
			break
		}

		if isLagged, isReplicating, err = s.checkSlaveStatus(); err == nil {
			return
		}

		time.Sleep(time.Second * 3)
		retry--
	}

	return
}

func (s *SQLRunner) checkSlaveStatus() (isLagged, isReplicating corev1.ConditionStatus, err error) {
	var rows *sql.Rows
	isLagged, isReplicating = corev1.ConditionUnknown, corev1.ConditionUnknown
	rows, err = s.db.Query("show slave status;")
	if err != nil {
		return
	}

	defer rows.Close()

	if !rows.Next() {
		if err = rows.Err(); err != nil {
			return
		}
		return corev1.ConditionFalse, corev1.ConditionFalse, nil
	}

	var cols []string
	cols, err = rows.Columns()
	if err != nil {
		return
	}

	scanArgs := make([]interface{}, len(cols))
	for i := range scanArgs {
		scanArgs[i] = &sql.RawBytes{}
	}

	err = rows.Scan(scanArgs...)
	if err != nil {
		return
	}

	slaveIOState := strings.ToLower(columnValue(scanArgs, cols, "Slave_IO_State"))
	slaveSQLRunning := columnValue(scanArgs, cols, "Slave_SQL_Running")
	lastSQLError := columnValue(scanArgs, cols, "Last_SQL_Error")
	secondsBehindMaster := columnValue(scanArgs, cols, "Seconds_Behind_Master")

	if stringInArray(slaveIOState, errorConnectionStates) {
		return isLagged, corev1.ConditionFalse, fmt.Errorf("Slave_IO_State: %s", slaveIOState)
	}

	if slaveSQLRunning != "Yes" {
		return isLagged, corev1.ConditionFalse, fmt.Errorf("Last_SQL_Error: %s", lastSQLError)
	}

	isReplicating = corev1.ConditionTrue

	var longQueryTime float64
	if err = s.GetGlobalVariable("long_query_time", &longQueryTime); err != nil {
		return
	}

	sec, _ := strconv.ParseFloat(secondsBehindMaster, 64)
	if sec > longQueryTime*100 {
		isLagged = corev1.ConditionTrue
	} else {
		isLagged = corev1.ConditionFalse
	}

	return
}

func (s *SQLRunner) CheckReadOnly() (corev1.ConditionStatus, error) {
	var readOnly uint8
	if err := s.GetGlobalVariable("read_only", &readOnly); err != nil {
		return corev1.ConditionUnknown, err
	}

	if readOnly == 0 {
		return corev1.ConditionFalse, nil
	}

	return corev1.ConditionTrue, nil
}

func (sr *SQLRunner) GetGlobalVariable(param string, val interface{}) error {
	query := fmt.Sprintf("select @@global.%s", param)
	return sr.db.QueryRow(query).Scan(val)
}

func (sr *SQLRunner) Close() error {
	return sr.db.Close()
}

func columnValue(scanArgs []interface{}, slaveCols []string, colName string) string {
	columnIndex := -1
	for idx := range slaveCols {
		if slaveCols[idx] == colName {
			columnIndex = idx
			break
		}
	}

	if columnIndex == -1 {
		return ""
	}

	return string(*scanArgs[columnIndex].(*sql.RawBytes))
}

func stringInArray(str string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, str)
	return index < len(strArray) && strArray[index] == str
}
