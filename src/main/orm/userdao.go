package orm

import (
	"database/sql"
	"main/pojo"
)

type UserInfoDao struct {
	sql.DB
}

/*
插入sql，使用preparedstatement避免sql注入
返回影响行数，err
*/
func (db *UserInfoDao) AddUserInfo(user *pojo.UserInfo) (ret int, err error) {
	stat, perr := db.Prepare("insert into user_info(id,name,age,mobile,email) values(?,?,?,?,?)")
	if perr != nil {
		return 0, perr
	}
	res, rerr := stat.Exec(user.Id, user.Name, user.Age, user.Mobile, user.Email)
	if rerr != nil {
		return 0, rerr
	}
	defer stat.Close()
	c, cerr := res.RowsAffected()
	return int(c), cerr
}

/*
更改，使用preparedstatement 返回影响行数，err
*/
func (db *UserInfoDao) UpdateUserInfoByIdAndMobile(user *pojo.UserInfo) (ret int, err error) {
	stat, perr := db.Prepare("update user_info set age=?,name=?,email=? where id=? and mobile=?")
	if perr != nil {
		return 0, perr
	}
	res, rerr := stat.Exec(user.Age, user.Name, user.Email, user.Id, user.Mobile)
	if rerr != nil {
		return 0, rerr
	}
	defer stat.Close()
	c, cerr := res.RowsAffected()
	return int(c), cerr
}

/*
select 列表结果集转换  将rows转换为对象列表，没有使用反射，使用对象new再赋值
*/
func parseUserFromRows(rows *sql.Rows, initCap int) (users []pojo.UserInfo, err error) {
	var id int32
	var name string
	var email string
	var age int8
	var mobile string
	ret := make([]pojo.UserInfo, initCap)
	num := 0
	for rows.Next() {
		//scan 查询的顺序必须和select字段顺序一致，不然会出现赋值错误
		serr := rows.Scan(&id, &name, &age, &mobile, &email)
		if serr != nil {
			us := make([]pojo.UserInfo, 0)
			return us, serr
		} else {
			us := pojo.UserInfo{id, name, age, mobile, email,nil}
			if num > initCap {
				ret = append(ret, us)
				num++
			} else {
				ret[num] = us
				num++
			}

		}
	}
	return ret[:num], nil
}

func (db *UserInfoDao) GetUserInfoById(id int32) (user *pojo.UserInfo, err error) {
	stat, perr := db.Prepare("select id,name,age,mobile,email from user_info where id=?")
	if perr != nil {
		return nil, perr
	}
	rows, rerr := stat.Query(id)
	if rerr != nil {
		return nil, rerr
	}
	users, uerr := parseUserFromRows(rows, 1)
	defer stat.Close()
	if uerr != nil {
		return nil, uerr
	} else {
		if len(users) > 0 {
			return &users[0], nil
		} else {
			return nil, nil
		}
	}
}

/*
通过email 拿到列表，返回对象列表和err
*/
func (db *UserInfoDao) GetUserInfoListByEmail(email string) (users []pojo.UserInfo, err error) {
	stat, perr := db.Prepare("select id,name,age,mobile,email from user_info where email=?")
	if perr != nil {
		return nil, perr
	}
	rows, rerr := stat.Query(email)
	if rerr != nil {
		return nil, rerr
	}
	us, uerr := parseUserFromRows(rows, 20)
	defer stat.Close()
	if uerr != nil {
		return nil, uerr
	} else {
		return us, nil
	}
}

/*
使用sql preparedstatment 执行，返回影响行数，err
*/
func (db *UserInfoDao) DeleteById(id int32) (c int, err error) {
	var deleteByIdSQL = "delete from user_info where id=?"
	stat, perr := db.Prepare(deleteByIdSQL)
	if perr != nil {
		return 0, perr
	}
	res, rerr := stat.Exec(id)
	if rerr != nil {
		return 0, rerr
	}
	defer stat.Close()
	cc, cerr := res.RowsAffected()
	return int(cc), cerr
}
