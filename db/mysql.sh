#!/bin/bash
HOSTNAME="127.0.0.1"                                           #数据库信息
PORT="3306"
USERNAME="root"
PASSWORD="123456"
DBNAME="registerDB"                                                #数据库名称
TABLENAME="userTbl"                                            #数据库中表的名称
 
#创建数据库
create_db_sql="create database IF NOT EXISTS ${DBNAME} CHARACTER SET 'utf8' COLLATE 'utf8_general_ci'"
mysql -h${HOSTNAME}  -P${PORT}  -u${USERNAME} -p${PASSWORD} -e "${create_db_sql}"
 
#create program table
create_user_table="create table IF NOT EXISTS ${TABLENAME} (
#userID bigint UNSIGNED NOT NULL Primary KEY AUTO_INCREMENT,
name varchar(64),
password varchar(128) default '',
mobile varchar(11) Primary KEY,  #key
gender varchar(1),
birthday date)"
mysql -h${HOSTNAME}  -P${PORT}  -u${USERNAME} -p${PASSWORD} ${DBNAME} -e "${create_user_table}"
