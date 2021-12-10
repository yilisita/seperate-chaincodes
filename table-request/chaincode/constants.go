// 每次只能打开一个工程文件夹
package chaincode

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (

	// 通道名
	MYCHANNEL = "mychannel"

	// 链码:table
	TABLE = "table"

	// 链码:account
	ACCOUNT = "account"

	// 链码:table-request
	TABLE_REQUEST = "table-request"

	// 链码:attribute-request
	ATTRIBUTE_REQUEST = "attribute-request"

	// 私有数据库: 是否需要?
	//requestCollection = "requestCollection"

	// ORG1的私有数据库, 所属链码:table
	myCollection = "Org1PrivateCollection"

	// ORG2的表格数据库, 所属链码：table-request
	// 想想是否需要单独的tableCollection
	tableCollection = "tableCollection"

	// ORG2的购买结果私有数据库，所属链码：attribute-request
	org2Collection = "Org2Collection"

	// 交易属性一次需要花费的金钱, 100元
	ATTRIBUTE_FEE float64 = 100

	// 交易表格一次需要花费的金钱，1000元
	TABLE_FEE float64 = 1000
)

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

var ATTRIBUTES_STRING = []string{
	"装机容量",
	"本月发电量",
	"本季度发电量",
	"截至本月累计发电量",
	"本月上网电量",
	"本季度上网电量",
	"截至本月上网电量",
	"综合厂本月用电量",
	"综合厂本季度用电量",
	"综合长截至本月用电量",
	"本月自发用电量",
	"本季度自发用电量",
	"截至本月自发用电量",
}

var PROPOSALS_STRINGS = []string{"求和", "求平均值", "计算相关系数"}

var PROPOSALS = []string{"GetAttributeTotal", "GetAttributeAve", "GetPearson"}

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

// 数据需求方发送请求的结构
/*--- 七个字段 ----*/
type Request struct {
	ID       string `json:"ID"`
	Demander string `json:"Demander"`
	// 目标表格
	TableID string `json:"TableID"`
	// 目标字段
	AttributeID []int `json:"Attribute"`
	// 计算内容 : 求和 or 平均值
	Proposal    int    `json:"Proposal"` // 0, 1, 2 (求和,平均值，相关系数)
	RequestTime string `json:"RequesTime"`
	Complete    bool   `json:"Complete"`
}

// 数据提供方返回响应的结构
/*---- 六个字段 ----*/
type Response struct {
	ID          string   `json:"ID"`
	TableName   string   `json:"TableName"`
	Attribute   []string `json:"Attribute"`
	ProposalStr string   `json:"Target"`
	RequestTime string   `json:"RequestTime"`
	// 结果
	Amount float64 `json:"Amount"`
}
