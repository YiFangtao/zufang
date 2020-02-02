package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"zufang/Ihomeweb/pkg/logging"
	"zufang/Ihomeweb/utils"
)

type User struct {
	Id           int           `json:"user_id"`                    //用户编号
	Name         string        `orm:"size(32)" json:"name"`        //用户昵称
	PasswordHash string        `orm:"size(128)" json:"password"`   //用户加密后的密码
	Mobile       string        `orm:"size(11)" json:"mobile"`      //手机号
	RealName     string        `orm:"size(32)" json:"real_name"`   //真实姓名
	IdCard       string        `orm:"size(20)" json:"id_card"`     //身份证号
	AvatarUrl    string        `orm:"size(256)" json:"avatar_url"` //用户头像路径
	Houses       []*House      `orm:"reverse(many)" json:"houses"` //用户发布的房屋信息
	Orders       []*OrderHouse `orm:"reverse(many)" json:"orders"` //用户下的订单
}

//房屋信息
type House struct {
	Id            int           `json:"house_id"`                                          //房屋编号
	User          *User         `orm:"rel(fk)" json:"user_id"`                             //房屋主人的用户编号 与用户进行关联
	Area          *Area         `orm:"rel(fk)" json:"area_id"`                             //所属地的区域编号
	Title         string        `orm:"size(64)" json:"title"`                              //房屋标题
	Price         int           `orm:"default(0)" json:"price"`                            //单价 单位：分 每次的价格都乘以100
	Address       string        `orm:"size(512)" orm:"default('')" json:"address"`         //地址
	RoomCount     int           `orm:"default(1)" json:"room_count"`                       //房间数
	Acreage       int           `orm:"default(0)" json:"acreage"`                          //房间总面积
	Unit          string        `orm:"size(32)" orm:"default('')" json:"unit"`             //房屋单位，如几室几厅
	Capacity      int           `orm:"default(1)" json:"capacity"`                         //房屋容纳的总人数
	Beds          string        `orm:"size(64)" orm:"default('')" json:"beds"`             //房屋床铺的配置
	Deposit       int           `orm:"default(0)" json:"deposit"`                          //押金
	MinDays       int           `orm:"default(1)" json:"min_days"`                         //最少入住的天数
	MaxDays       int           `orm:"default(0)" json:"max_days"`                         //最多入住的天数 0表示不限制
	OrderCount    int           `orm:"default(0)" json:"order_count"`                      //预计完成的该房屋的订单数
	IndexImageUrl string        `orm:"size(256)" orm:"default('')" json:"index_image_url"` //房屋主图片
	Facilities    []*Facility   `orm:"reverse(many)" json:"facilities"`                    //房屋设施
	Images        []*HouseImage `orm:"reverse(many)" json:"img_urls"`                      //房屋图片
	Orders        []*OrderHouse `orm:"reverse(many)" json:"orders"`                        //房屋的订单
	Ctime         time.Time     `orm:"auto_now_add;type(datetime)" json:"ctime"`
}

//房屋订单信息
type OrderHouse struct {
	Id         int       `json:"order_id"`                            //订单编号
	User       *User     `orm:"rel(fk)" json:"user_id"`               //下单的用户编号
	House      *House    `orm:"rel(fk)" json:"house_id"`              //预定的房间编号
	BeginDate  time.Time `orm:"type(datetime)" json:"begin_date"`     //预定的起始时间
	EndDate    time.Time `orm:"type(datetime)" json:"end_date"`       //预定的结束时间
	Days       int       `json:"days"`                                //预定的总天数
	HousePrice int       `json:"house_price"`                         //房屋的单价
	Amount     int       `json:"amount"`                              //订单总金额
	Status     string    `orm:"default(WAIT_ACCEPT)" json:"status"`   //订单状态
	Comment    string    `orm:"size(512)" json:"comment"`             //订单评论
	Ctime      time.Time `orm:"auto_now;type(datetime)" json:"ctime"` //每次更新此表都会更新这个字段
	Credit     bool      `json:"credit"`                              //表示个人征信 true表示良好
}

type Area struct {
	Id     int      `json:"area_id"`                    //区域编号
	Name   string   `orm:"size(32)" json:"area_name"`   //区域名字
	Houses []*House `orm:"reverse(many)" json:"houses"` //区域所有的房屋 与房屋表进行关联
}

type Facility struct {
	Id     int      `json:"facility_id"`                  //设施编号
	Name   string   `orm:"size(32)" json:"facility_name"` //设施名字
	Houses []*House `orm:"rel(m2m)" json:"houses"`        //哪些房屋有此设施
}

type HouseImage struct {
	Id    int    `json:"house_image_id"`         //图片id
	Url   string `orm:"size(256)" json:"url"`    //图片url
	House *House `orm:"rel(fk)" json:"house_id"` //图片所属房屋编号
}

var (
	HOME_PAGE_MAX_HOUSES    = 5 //首页最高展示的房屋数量
	HOME_LIST_PAGE_CAPACITY = 2 //房屋列表页面每页显示条目数
)

const (
	ORDER_STATUS_WAIT_ACCEPT  = "WAIT_ACCEPT"  //待接单
	ORDER_STATUS_WAIT_PAYMENT = "WAIT_PAYMENT" //待支付
	ORDER_STATUS_PAID         = "PAID"         //已支付
	ORDER_STATUS_WAIT_COMMENT = "WAIT_COMMENT" //待评论
	ORDER_STATUS_COMPLETE     = "COMPLETE"     //已完成
	ORDER_STATUS_CANCELED     = "CANCELED"     //已取消
	ORDER_STATUS_REJECTED     = "REJECTED"     //已拒单
)

//处理房子信息
func (this *House) ToHouseInfo() interface{} {
	houseInfo := map[string]interface{}{
		"house_id":    this.Id,
		"title":       this.Title,
		"price":       this.Price,
		"area_name":   this.Area.Name,
		"img_url":     utils.AddDomain2Url(this.IndexImageUrl),
		"room_count":  this.RoomCount,
		"order_count": this.OrderCount,
		"address":     this.Address,
		"user_avatar": utils.AddDomain2Url(this.User.AvatarUrl),
		"ctime":       this.Ctime.Format("2006-01-02 15:04:05"),
	}
	return houseInfo
}

//处理1个房子的全部信息
func (this *House) ToOneHouseDesc() interface{} {
	houseDesc := map[string]interface{}{
		"house_id":    this.Id,
		"user_id":     this.User.Id,
		"user_name":   this.User.Name,
		"user_avatar": utils.AddDomain2Url(this.User.AvatarUrl),
		"title":       this.Title,
		"price":       this.Price,
		"address":     this.Address,
		"room_count":  this.RoomCount,
		"acreage":     this.Acreage,
		"unit":        this.Unit,
		"capacity":    this.Capacity,
		"beds":        this.Beds,
		"deposit":     this.Deposit,
		"min_days":    this.MinDays,
		"max_days":    this.MaxDays,
	}
	//房屋图片
	var imaUrls []string
	for _, imaUrl := range this.Images {
		imaUrls = append(imaUrls, utils.AddDomain2Url(imaUrl.Url))
	}
	houseDesc["img_urls"] = imaUrls

	//房屋设施
	var facilities []int
	for _, facility := range this.Facilities {
		facilities = append(facilities, facility.Id)
	}
	houseDesc["facilities"] = facilities

	//评论信息
	var comments []interface{}
	var orders []OrderHouse
	o := orm.NewOrm()
	orderNum, err := o.QueryTable("order_house").Filter("house_id", this.Id).Filter("status", ORDER_STATUS_COMPLETE).OrderBy("ctime").Limit(10).All(&orders)
	if err != nil {
		logging.Debug("查询订单信息失败", err)
	}
	for i := 0; i < int(orderNum); i++ {
		o.LoadRelated(&orders[i], "User")
		var username string
		if orders[i].User.Name == "" {
			username = "匿名用户"
		} else {
			username = orders[i].User.Name
		}
		comment := map[string]string{
			"comment":   orders[i].Comment,
			"user_name": username,
			"ctime":     orders[i].Ctime.Format("2006-01-02 15:04:05"),
		}
		comments = append(comments, comment)
	}
	houseDesc["comments"] = comments

	return houseDesc
}

//处理订单信息
func (this *OrderHouse) ToOrderInfo() interface{} {
	orderInfo := map[string]interface{}{
		"order_id":   this.Id,
		"title":      this.House.Title,
		"img_url":    utils.AddDomain2Url(this.House.IndexImageUrl),
		"start_date": this.BeginDate.Format("2006-01-02 15:04:05"),
		"end_date":   this.EndDate.Format("2006-01-02 15:04:05"),
		"ctime":      this.Ctime.Format("2006-01-02 15:04:05"),
		"days":       this.Days,
		"amount":     this.Amount,
		"status":     this.Status,
		"comment":    this.Comment,
		"credit":     this.Credit,
	}
	return orderInfo
}

func init() {
	//注册mysql驱动
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//设置默认数据库
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp("+utils.G_mysql_addr+":"+utils.G_mysql_port+")/gomicro?charset=utf8", 30)
	//orm.RegisterDataBase("default", "mysql", "root:123456@tcp(192.168.3.131:3306)/gomicro?charset=utf8", 30)
	//注册model
	orm.RegisterModel(new(User), new(House), new(Area), new(Facility), new(HouseImage), new(OrderHouse))

	orm.RunSyncdb("default", false, true)
}
