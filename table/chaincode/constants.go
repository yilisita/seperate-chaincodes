// 每次只能打开一个工程文件夹
package chaincode

import (
	"table/paillier"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const requestCollection = "requestCollection"

// ORG1的私有数据库
const myCollection = "Org1PrivateCollection"

// ORG的表格数据库
const tableCollection = "tableCollection"

var ATTRIBUTES = []string{
	// 装机容量
	"ZhuangJiAmount",
	// 本月发电量
	"GenAmountThisMonth",
	// 本季累计发电量
	"GenAmountThisSeason",
	// 截止本月发电量
	"GenAmountTillThisMonth",
	// 本月上网电量
	"ShangWangAmountThisMonth",
	// 本季累计上网电量
	"ShangWangAmountThisSeason",
	// 本月止累计
	"ShangWangAmountTillThisMonth",
	// 综合厂本月用电量
	"ZongheThisMonth",
	// 综合厂本季度用电量
	"ZongheThisSeason",
	// 综合厂截止本月用电量
	"ZongheTillThisMonth",
	// 自发本月用电量
	"SelfGenThisMonth",
	// 自发本季度用电量
	"SelfGenThisSeason",
	// 自发截止本月用电量
	"SelfGenTillThisMonth",
}

type SmartContract struct {
	contractapi.Contract
	Tablerequest []BuyTable
}

// 在逻辑上，购买字段的请求和购买整个表的请求应该是分开的
type BuyTable struct {
	ID       string `json:"ID"`
	Demander string `json:"Demander"`
	// 目标表格
	TableID     string `json:"TableID"`
	RequestTime string `json:"RequesTime"`
	Complete    bool   `json:"Complete"`
}

// 使用paillier进行加密。确保上传到链上的数据是整数
var p = paillier.KeyGenPaillier()

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism accross languages
// golang keeps the order when marshal to json but doesn't order automatically
type Asset struct {
	// 用户名
	UserName string `json:"UserName"`
	// 用户ID
	UserID string `json:"UserID"`
	// 装机容量
	ZhuangJiAmount int64 `json:"ZhuangJiAmount"`
	// 本月发电量
	GenAmountThisMonth int64 `json:"GenAmountThisMonth"`
	// 本季累计发电量
	GenAmountThisSeason int64 `json:"GenAmountThisSeason"`
	// 截止本月发电量
	GenAmountTillThisMonth int64 `json:"GenAmountTillThisMonth"`
	// 本月上网电量
	ShangWangAmountThisMonth int64 `json:"ShangWangAmountThisMonth"`
	// 本季累计上网电量
	ShangWangAmountThisSeason int64 `json:"ShangWangAmountThisSeason"`
	// 本月止累计
	ShangWangAmountTillThisMonth int64 `json:"ShangWangAmountTillThisMonth"`
	// 综合厂本月用电量
	ZongheThisMonth int64 `json:"ZongheThisMonth"`
	// 综合厂本季度用电量
	ZongheThisSeason int64 `json:"ZongheThisSeason"`
	// 综合厂截止本月用电量
	ZongheTillThisMonth int64 `json:"ZongheTillThisMonth"`
	// 自发本月用电量
	SelfGenThisMonth int64 `json:"SelfGenThisMonth"`
	// 自发本季度用电量
	SelfGenThisSeason int64 `json:"SelfGenThisSeason"`
	// 自发截止本月用电量
	SelfGenTillThisMonth int64 `json:"SelfGenTillThisMonth"`
}

// Encrypted用来描述加密字段,flags中有13个bool值，true：表示该索引对应的字段被加密了。0.
type Table struct {
	TableName string  `json:"TableName"`
	TableID   string  `json:"TableID"`
	Assets    []Asset `json:"Assets"`
	Flags     []int   `json:"Flags"`
}

// TableView存在于ledger中，所有人可见，用来描述当前有哪些表格;
type TableView struct {
	TableName string `json:"TableName"`
	TableID   string `json:"TableID"`
}
