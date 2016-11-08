package index

import (
	"database/sql"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

func StarZtreeUpdateIncrTask(endpoint string, conn *sql.DB) {
	err := ztree(endpoint, conn)
	if err != nil {
		log.Println(err)
	}
}

func ztree(endpoint string, conn *sql.DB) error {
	grp := strings.Split(endpoint, ".")
	len := len(grp)
	if len < 0 {
		return nil
	}

	if m := net.ParseIP(endpoint); m != nil {
		return nil
	}
	var group, idc, service string
	var grp_id, idc_id, service_id string
	if len > 3 {
		group = grp[len-1]
		idc = grp[len-2]
		service = strings.Join(grp[1:len-2], ".")
		err := conn.QueryRow("SELECT id FROM ztree WHERE name = ?", group).Scan(&grp_id)
		if err != nil {
			grp_id, err = add("0", group, conn)
			if err != nil {
				return err
			}
		}

		err = conn.QueryRow("SELECT id FROM ztree WHERE name = ? and pid = ?", idc, grp_id).Scan(&idc_id)
		if err != nil {
			idc_id, err = add(grp_id, idc, conn)
			if err != nil {
				return err
			}
		}

		err = conn.QueryRow("SELECT id FROM ztree WHERE name = ? and pid = ?", service, idc_id).Scan(&service_id)
		if err != nil {
			service_id, err = add(idc_id, service, conn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isip(ip string) (b bool) {
	if m, _ := regexp.MatchString("^[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}$", ip); !m {
		return false
	}
	return true
}

func add(pid, name string, conn *sql.DB) (string, error) {
	var last_id string
	var id string
	err := conn.QueryRow("SELECT id FROM ztree WHERE pid = ? ORDER BY id DESC LIMIT 1", pid).Scan(&last_id)
	if err != nil {
		switch len(pid) {
		case 1:
			id = "10"
		case 2:
			id = pid + "10"
		case 4:
			id = pid + "1000"
		default:
			return id, err
		}
	} else {
		tmp_id, err := strconv.ParseUint(last_id, 10, 64)
		if err != nil {
			return id, err
		}
		tmp_id++
		id = strconv.FormatUint(tmp_id, 10)
	}

	sqlStr := "INSERT INTO ztree VALUES(?,?,?,?)"
	_, err = conn.Exec(sqlStr, id, pid, name, "true")
	if err != nil {
		return id, err
	}
	log.Println(id, pid, name)
	return id, nil
}
