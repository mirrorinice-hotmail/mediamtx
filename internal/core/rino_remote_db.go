package core

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var gRNRemoteDbMgr TRNRemoteDbMgr

////////////////////////////////////

// const (
// 	test_host      = "localhost"
// 	test_port      = 5432
// 	test_user      = "postgres"
// 	test_pw        = "rino1234"
// 	test_dbname    = "test_rino_cctv_list"
// 	test_tablename = "tbl_cctv_info"
// )

const (
	col_mgr_no      = "mgr_no"
	col_cctv_ip     = "ip_addr"
	col_port_num    = "port_num"
	col_stream_path = "stream_path"
	col_stream_pw   = "stream_pw"
	col_cctv_nm     = "cctv_nm"
	col_addr1       = "addr1"
	col_addr2       = "addr2"
	col_serial_num  = "serial_num"
	col_manager_nm  = "manager_nm"
	col_stream_id   = "stream_id"
	col_rtsp_01     = "rtsp_01"
	col_rtsp_02     = "rtsp_02"
	col_health      = "health"
)

type TRNRemoteDbMgr struct {
	Name     string
	DbmsInfo DbmsST
}

func (obj *TRNRemoteDbMgr) init(dbmsInfo *DbmsST) {
	obj.Name = "CctvListMgr"
	obj.DbmsInfo = *dbmsInfo
}

func (obj *TRNRemoteDbMgr) db_open() *sql.DB {
	psinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", obj.DbmsInfo.Host, obj.DbmsInfo.Port, obj.DbmsInfo.User, obj.DbmsInfo.Pass, obj.DbmsInfo.Dbname)

	db, err := sql.Open("postgres", psinfo)
	if err != nil {
		fmt.Println(obj.Name, ": DB connection failure")
	} else {
		fmt.Println(obj.Name, ": DB connection success")
	}

	return db
}

func (obj *TRNRemoteDbMgr) Read_stream_list() TRNStreamsMAP {

	var newStream TRNStreamsMAP
	remote_db := obj.db_open()
	if remote_db == nil {
		return newStream
	}
	defer remote_db.Close()

	query := "SELECT " +
		col_stream_id + "," +
		col_rtsp_01 + "," +
		col_rtsp_02 + "," +
		col_cctv_nm + "," +
		col_cctv_ip + "," +
		col_health +
		" FROM " + obj.DbmsInfo.TableName
	fmt.Printf("sql : query(%s)\n", query)
	rows, err := remote_db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	newStream = makeTemporalStreams(rows)
	return newStream
}

// ///////////////////////////////////////////////////////////////////////////////

func makeTemporalStreams(rows *sql.Rows) TRNStreamsMAP {

	var newStreamsList = make(TRNStreamsMAP)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("_makeTemporalStreams", ": recovered from panic:", r)
		}
	}()

	// query := "SELECT " +
	// 	col_stream_id + "," +
	// 	col_rtsp_01 + "," +
	// 	col_rtsp_02 + "," +
	// 	col_cctv_nm + "," +
	// 	col_cctv_ip + "," +
	// 	col_health +
	for rows.Next() {
		var val_stream_id, val_rtsp_01, val_rtsp_02, val_cctv_nm, val_cctv_ip string
		var val_health byte
		err := rows.Scan(&val_stream_id, &val_rtsp_01, &val_rtsp_02, &val_cctv_nm, &val_cctv_ip, &val_health)
		if err != nil {
			fmt.Printf("Warning! can't parse stream in stream_id(%s)\n", val_stream_id)
			continue
		}
		fmt.Printf("stream list: stream_id(%s), rtsp_01(%s), cctv_nm(%s)\n",
			val_stream_id, val_rtsp_01, val_cctv_nm)

		tmpStream := TRNStreamST{
			//StreamId: val_stream_id,
			//CctvName: val_cctv_nm,
			//CctvIp:   val_cctv_ip,
			//Channels:     make(ChannelMAP),
			RtspUrl: val_rtsp_01,
			//RtspUrl_2: val_rtsp_02,
			//Status:    val_health,
			OnDemand: true,
			//DisableAudio: true,
			//Debug:        false,
		}
		//tmpStream.Channels["0"] = ChannelST{}
		newStreamsList[val_stream_id] = tmpStream
	}
	return newStreamsList
}

func (obj *TRNRemoteDbMgr) db_add_samples() bool {
	var insert_val_list = []string{
		"VALUES ('cctv_17_4_1' , '10.17.4.51' , '5432' , '용주면 가호길 91' , 'cctv_17_4_1' , 'rtsp://admin:hap_1000!%40%23@10.17.4.50:558/LiveChannel/00/media.smp' , 'rtsp://hcmanager:hap_1000!%40%23@10.17.4.51:554/Profile02/media.smp' )",
		"VALUES ('cctv_17_5_1' , '10.17.5.51' , '5432' , '용주면 고품3길 27' , 'cctv_17_5_1' , 'rtsp://admin:hap_1000!%40%23@10.17.5.50:558/LiveChannel/00/media.smp' , 'rtsp://hcmanager:hap_1000!%40%23@10.17.5.51:554/Profile02/media.smp' )",
	}
	remote_db := obj.db_open()
	if remote_db == nil {
		return false
	}
	defer remote_db.Close()

	const query_base = "INSERT INTO tbl_cctv_info (mgr_no, ip_addr, port_num, cctv_nm, stream_id, rtsp_02, rtsp_01) "

	for i, val := range insert_val_list {
		query := query_base + val
		fmt.Printf("%d : query(%s)\n", i, query)
		_, err := remote_db.Query(query)
		if err != nil {
			panic(err) //->return
		}
	}
	return true
}

/*
func (obj *tCctvListMgr) db_add(db *sql.DB, in_no string) bool {
	const name = "cctvlist_mgr"
	if in_no == "" {
		return false
	}
	query_action := "INSERT INTO"
	val_mgr_no := in_no
	val_cctv_ip := "10.10.0." + val_mgr_no
	val_port_num := "5432"
	val_stream_path := "rtsp://10.10.0.12:5564/" + val_mgr_no + "/stream01"
	val_cctv_nm := "CCTV_" + val_mgr_no
	val_serial_num := "10010" + val_mgr_no
	val_stream_id := "cctv002"
	val_rtsp_01 := "rtsp://210.99.70.120:1935/live/cctv001.stream"
	val_rtsp_02 := "rtsp://"

	query := query_action + " " + obj.DbmsInfo.TableName +
		" (" +
		col_mgr_no + ", " +
		col_cctv_ip + ", " +
		col_port_num + ", " +
		col_cctv_nm + ", " +
		col_stream_path + ", " +
		col_serial_num + ", " +
		col_stream_id + ", " +
		col_rtsp_01 + ", " +
		col_rtsp_02 +
		") " +
		" VALUES ($1, $2, $3, $4, $5, $6 , $7, $8, $9) ;"

	log.Println(name, ": db add(", query, ")")
	_, err := db.Exec(query,
		val_mgr_no,
		val_cctv_ip,
		val_port_num,
		val_cctv_nm,
		val_stream_path,
		val_serial_num,
		val_stream_id,
		val_rtsp_01,
		val_rtsp_02,
	)

	return obj.db_result_print(err, query_action)
}

func (obj *tCctvListMgr) db_update(db *sql.DB, in_no string) bool {
	const name = "cctvlist_mgr"
	if in_no == "" {
		return false
	}

	query_action := "UPDATE"
	val_mgr_no := in_no
	val_cctv_nm := "CCTV_" + val_mgr_no
	query := query_action + " " + obj.DbmsInfo.TableName +
		" SET " + col_cctv_nm + " = $1" +
		" WHERE " + col_mgr_no + " = $2 ;"
	_, err := db.Exec(query,
		val_cctv_nm,
		val_mgr_no)

	return obj.db_result_print(err, query_action)
}

func (obj *tCctvListMgr) db_delete(db *sql.DB, in_no string) bool {
	const name = "cctvlist_mgr"
	if in_no == "" {
		return false
	}
	query_action := "DELETE FROM"
	val_mgr_no := in_no
	query := query_action + " " + obj.DbmsInfo.TableName +
		" WHERE " + col_mgr_no + " = $1;"
	_, err := db.Exec(query,
		val_mgr_no)

	return obj.db_result_print(err, query_action)
}

func (obj *tCctvListMgr) db_read(db *sql.DB, in_no string) bool {
	const name = "cctvlist_mgr"
	if in_no == "" {
		return false
	}

	query_action := "SELECT"

	query := query_action + " " +
		col_mgr_no + "," + col_stream_path + "," + col_cctv_nm +
		" FROM " + obj.DbmsInfo.TableName
	if in_no != "all" {
		val_mgr_no := in_no
		query = query +
			" WHERE " + col_mgr_no + " = '" + val_mgr_no + "' ;"
	}
	fmt.Printf("%s : query(%s)\n", name, query)
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var val_mgr_no string
		var val_stream_path string
		var val_cctv_num string
		err := rows.Scan(&val_mgr_no, &val_stream_path, &val_cctv_num)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s : mgr_no(%s), path(%s), name(%s)\n", name, val_mgr_no, val_stream_path, val_cctv_num)
	}

	return obj.db_result_print(err, query_action)
}

func (obj *tCctvListMgr) db_result_print(err error, in_queryaction string) bool {
	const name = "cctvlist_mgr"
	if err != nil {
		fmt.Println(name, ": err(", err.Error(), ")")
		return false
	} else {
		fmt.Println(name, ": ", in_queryaction, ` success!`)
		return true
	}
}
*/
