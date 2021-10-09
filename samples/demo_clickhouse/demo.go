package demo_clickhouse

import (
	"fmt"
	"github.com/rubinus/zgo"
	"log"
	"time"
)

/*
@Time : 2019-09-26 16:34
@Author : rubinus.chu
@File : demo
@project: origin
*/

const (
	label_bj = "clickhouse"
	//label_bj = "1934193167338"
	//label_sh = "1049232219283"
)

func Get() {
	connChan, err := zgo.CK.GetConnChan()
	if err != nil {
		zgo.Log.Error(err)
		return
	}
	if connect, ok := <-connChan; ok {
		_, err = connect.Exec(`
		CREATE TABLE IF NOT EXISTS example (
			country_code FixedString(2),
			os_id        UInt8,
			browser_id   UInt8,
			categories   Array(Int16),
			action_day   Date,
			action_time  DateTime
		) engine=Memory
	`)

		if err != nil {
			log.Fatal(err)
		}

		var (
			tx, _   = connect.Begin()
			stmt, _ = tx.Prepare("INSERT INTO example (country_code, os_id, browser_id, categories, action_day, action_time) VALUES (?, ?, ?, ?, ?, ?)")
		)
		defer stmt.Close()

		for i := 0; i < 5; i++ {
			if _, err := stmt.Exec(
				"RU",
				10+i,
				100+i,
				//clickhouse.Array([]int16{1, 2, 3}),
				[]int16{1, 2, 3, 4, 5, 6},
				time.Now(),
				time.Now(),
			); err != nil {

				log.Fatal(err)
			}
		}

		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}

		rows, err := connect.Query("SELECT country_code, os_id, browser_id, categories, action_day, action_time FROM example")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var (
				country               string
				os, browser           uint8
				categories            []int16
				actionDay, actionTime time.Time
			)
			if err := rows.Scan(&country, &os, &browser, &categories, &actionDay, &actionTime); err != nil {
				log.Fatal(err)
			}
			log.Printf("country: %s, os: %d, browser: %d, categories: %v, action_day: %s, action_time: %s", country, os, browser, categories, actionDay, actionTime)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("no connection")
	}
}