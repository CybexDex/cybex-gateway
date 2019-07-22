package modeladmin

import (
	"cybex-gateway/types"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/shopspring/decimal"
)

// Asset ...
type Asset struct {
	gorm.Model
	Name         string `gorm:"unique" json:"name"`
	JadeName     string `json:"jadeName"`
	Blockchain   string `json:"blockchain"`
	Projectname  string `json:"projectname"`
	CYBName      string `json:"cybname"`
	CYBID        string `json:"cybid"`
	Confirmation string `json:"confirmation"`

	SmartContract   string `json:"smartContract"`
	GatewayAccount  string `json:"gatewayAccount"`           // 以account为准
	GatewayID       string `json:"-"`                        // 方便读块时排除不需要的部分
	GatewayPass     string `json:"-"`                        // 非常重要
	GatewayPassword string `gorm:"-" json:"gatewayPassword"` // 非常重要
	WithdrawPrefix  string `json:"withdrawPrefix"`

	DepositSwitch  *bool `json:"depositSwitch"`
	WithdrawSwitch *bool `json:"withdrawSwitch"`

	MinDeposit  decimal.Decimal `gorm:"type:numeric" json:"minDeposit"`
	MinWithdraw decimal.Decimal `gorm:"type:numeric" json:"minWithdraw"`
	WithdrawFee decimal.Decimal `gorm:"type:numeric" json:"withdrawFee"`
	DepositFee  decimal.Decimal `gorm:"type:numeric" json:"depositFee"`

	Precision string         `json:"precision"`
	ImgURL    string         `json:"imgURL"`
	HashLink  string         `json:"hashLink"`
	Info      postgres.Jsonb `gorm:"default:'{}'" json:"info"`
	UseMemo   *bool          `json:"useMemo"`
	Disabled  *bool          `json:"disabled"`
}

// ValidateCreate ...
func (a Asset) ValidateCreate() error {
	return validation.ValidateStruct(&a,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Name, validation.Required),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Blockchain, validation.Required),
		validation.Field(&a.CYBName, validation.Required),
		validation.Field(&a.GatewayAccount, validation.Required),
		validation.Field(&a.GatewayPassword, validation.Required),
		validation.Field(&a.WithdrawPrefix, validation.Required),
		validation.Field(&a.MinDeposit, validation.Required),
		validation.Field(&a.MinWithdraw, validation.Required),
		validation.Field(&a.WithdrawFee, validation.Required),
	)
}

// Save ...
func (j *Asset) Save() error {
	return db.Save(j).Error
}

// TestJP ...
func (j *Asset) TestJP() error {
	return nil
}

// TestCybex ...
func (j *Asset) TestCybex() error {
	return nil
}

const disableStr = "disabled is NULL or disabled is false"

// AssetsAll ...
func AssetsAll() (out []*Asset, err error) {
	err = db.Find(&out, disableStr).Error
	return out, err
}

// AssetsFind ...
func AssetsFind(asset string) (out *Asset, err error) {
	out = &Asset{}
	err = db.First(&out, &Asset{
		Name: asset,
	}, disableStr).Error
	return out, err
}

// AssetsCreate ...
func AssetsCreate(query *Asset) (*Asset, error) {
	err := db.Save(query).Error
	return query, err
}

// AssetsSwitch ...
func AssetsSwitch(query *types.Switch) (out []*Asset, err error) {
	// out = &Asset{}
	var names []string
	update := &Asset{}
	if query.Withdraw != nil {
		update.WithdrawSwitch = query.Withdraw
	}
	if query.Deposit != nil {
		update.DepositSwitch = query.Deposit
	}
	if query.Name != "" {
		names = strings.Split(query.Name, ",")
		err = db.Model(update).Where("name in (?)", names).UpdateColumns(update).Error
	} else {
		err = db.Model(update).UpdateColumns(update).Error
	}
	return out, err
}

// UpdateAsset ...
func UpdateAsset(query *Asset) (out *Asset, err error) {
	out = &Asset{}
	err = db.Find(out, query.ID).UpdateColumn(query).Error
	return out, err
}

// AssetsQuery ...
func AssetsQuery(query *Asset) (out []*Asset, err error) {
	err = db.Order("id desc").Find(&out, query).Error
	return out, err
}

// AssetsFrist ...
func AssetsFrist(query *Asset) (out *Asset, err error) {
	out = &Asset{}
	err = db.First(&out, query, disableStr).Error
	return out, err
}
