package main

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

func main() {
	err := dao()
	fmt.Printf("dao err:%+v", err)
}

func dao() error {
	return errors.Wrap(sql.ErrNoRows, "dao error")
}

/*
問題：
我們在數據庫操作的時候，比如 dao 層中當遇到一個 sql.ErrNoRows 的時候，是否應該 Wrap 這個 error，拋給上層。
為什麼，應該怎麼做請寫出代碼？


回答：
ErrNoRows產生於當QueryRow方法不返回行時，推遲到Scan方法返回ErrNoRows。
我們在遇到這個錯誤時，應當調用Wrap方法，捕獲錯誤的堆棧信息，即時定位錯誤，告知上層需要進行修改。
*/
