/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-ibax/packages/storage/sqldb"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"

	"github.com/shopspring/decimal"

	//"github.com/IBAX-io/go-explorer/conf"
	"strconv"
	"strings"

	"github.com/IBAX-io/go-ibax/packages/converter"
)

// Key is model
type Key struct {
	Ecosystem int64
	ID        int64           `gorm:"primary_key;not null"`
	PublicKey []byte          `gorm:"column:pub;not null"`
	Amount    decimal.Decimal `gorm:"not null"`
	Maxpay    decimal.Decimal `gorm:"not null"`
	Deposit   decimal.Decimal `gorm:"not null"`
	Multi     int64           `gorm:"not null"`
	Deleted   int64           `gorm:"not null"`
	Blocked   int64           `gorm:"not null"`
	Account   string          `gorm:"column:account;not null"`
	Lock      string          `gorm:"column:lock;type:jsonb"`
}

type KeyHex struct {
	Ecosystem     int64           `json:"ecosystem"`
	ID            string          `json:"id"`
	PublicKey     string          `json:"publickey"`
	Amount        string          `json:"amount"`
	Maxpay        decimal.Decimal `json:"maxpay"`
	Multi         int64           `json:"multi"`
	Deleted       int64           `json:"deleted"`
	Blocked       int64           `json:"blocked"`
	Ecosystemname string          `json:"ecosystemname"`
	TokenSymbol   string          `json:"token_symbol"`
}

type EcosyKeyHex struct {
	Ecosystem       int64  `json:"ecosystem"`
	IsValued        int64  `json:"isvalued"`
	Ecosystemname   string `json:"ecosystemname"`
	TokenSymbol     string `json:"token_symbol"`
	Amount          string `json:"amount"`
	Info            string `json:"info"`
	Emission_amount string `json:"emission_amount"`
	Type_emission   int64
	Type_withdraw   int64
}

type EcosyKeyTotalHex struct {
	Wallet         string          `json:"wallet"`
	Ecosystem      int64           `json:"ecosystem"`
	IsValued       int64           `json:"isvalued"`
	Ecosystemname  string          `json:"ecosystemname"`
	TokenSymbol    string          `json:"token_symbol"`
	Amount         string          `json:"amount"`
	Info           string          `json:"info"`
	EmissionAmount string          `json:"emission_amount"`
	MemberName     string          `json:"member_name"`
	MemberHash     string          `json:"member_hash"`
	LogoHash       string          `json:"logo_hash"`
	TypeEmission   int64           `json:"type_emission"`
	TypeWithdraw   int64           `json:"type_withdraw"`
	Transaction    int64           `json:"transaction"`
	InTx           int64           `json:"in_tx"`
	OutTx          int64           `json:"out_tx"`
	StakeAmount    decimal.Decimal `json:"stake_amount"`
	LockAmount     decimal.Decimal `json:"lock_amount"`
	InAmount       decimal.Decimal `json:"inamount"`
	OutAmount      decimal.Decimal `json:"outamount"`
	TxAmount       string          `json:"tx_amount"`
	TotalAmount    decimal.Decimal `json:"total_amount"`
	JoinTime       int64           `json:"join_time"`
	RolesName      string          `json:"roles_name"`
	Digits         int             `json:"digits"`
}

type EcosyKeyTotalDetail struct {
	Account     string          `json:"account"`
	Ecosystem   int64           `json:"ecosystem"`
	LogoHash    string          `json:"logo_hash"`
	Name        string          `json:"name"`
	JoinTime    int64           `json:"join_time"`
	TokenSymbol string          `json:"token_symbol"`
	Amount      string          `json:"amount"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	LockAmount  decimal.Decimal `json:"lock_amount"`
	StakeAmount decimal.Decimal `json:"stake_amount"`
	RolesName   string          `json:"roles_name"`
	Digits      int             `json:"digits"`
}

type AccountHistoryChart struct {
	Time        []int64  `json:"time"`
	Inamount    []string `json:"inamount"`
	Outamount   []string `json:"outamount"`
	TokenSymbol string   `json:"token_symbol"`
	Digits      int      `json:"digits"`
}

type KeysResult struct {
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
	Rets  []EcosyKeyTotalHex `json:"rets"`
}

type EcosyKeyList struct {
	Ecosystem    int64           `json:"ecosystem"`
	Account      string          `json:"account"`
	AccountName  string          `json:"account_name"`
	Amount       string          `json:"amount"`
	TotalAmount  decimal.Decimal `json:"total_amount"`
	AccountedFor decimal.Decimal `json:"accounted_for"`
	StakeAmount  decimal.Decimal `json:"stake_amount"`
	LockAmount   decimal.Decimal `json:"lock_amount"`
	TokenSymbol  string          `json:"token_symbol"`
	Digits       int             `json:"digits"`
}

type KeysListResult struct {
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Rets  []EcosyKeyList `json:"rets"`
}

type KeysListChartResult struct {
	KeysInfo KeysRet      `json:"keys_info"`
	KeyChart KeyInfoChart `json:"key_chart"`
}

type KeyLocks struct {
	NftMinerStake       string `json:"nft_miner_stake"`
	CandidateReferendum string `json:"candidate_referendum"`
	CandidateSubstitute string `json:"candidate_substitute"`
}

// SetTablePrefix is setting table prefix
func (m *Key) SetTablePrefix(prefix int64) *Key {
	m.Ecosystem = prefix
	return m
}

// TableName returns name of table
func (m Key) TableName() string {
	if m.Ecosystem == 0 {
		m.Ecosystem = 1
	}
	return `1_keys`
}

func (m *Key) Get(id int64, wallet string) (*EcosyKeyHex, error) {

	var (
		ecosystems []Ecosystem
	)
	da := EcosyKeyHex{}

	key := strconv.FormatInt(id, 10)
	wid, err := strconv.ParseInt(wallet, 10, 64)
	if err == nil {
		//
		err := conf.GetDbConn().Conn().Table("1_ecosystems").Where("id = ?", key).Find(&ecosystems).Error
		if err == nil {
			da.Ecosystem = id
			da.Ecosystemname = ecosystems[0].Name
			da.IsValued = ecosystems[0].IsValued
			if da.Ecosystem == 1 {
				da.TokenSymbol = SysTokenSymbol
			} else {
				da.TokenSymbol = ecosystems[0].TokenSymbol
			}
			//da.Token_title = ecosystems[0].Token_title
		}
		err = conf.GetDbConn().Conn().Table("1_keys").Where("id = ?", wid).Find(m).Error
		if err == nil {
			da.Ecosystem = id
			da.Amount = m.Amount.String()
		}
	}

	//err := conf.GetDbConn().Conn().Where("id = ? and ecosystem = ?", wallet, m.ecosystem).First(m).Error
	return &da, err
}

func (ts *Key) GetKeys(id int64, page int, size int, order string) (*[]KeyHex, int64, error) {
	var (
		tss    []Key
		ret    []KeyHex
		num    int64
		ioffet int
		i      int
	)

	if order == "" {
		order = "id asc"
	}
	num = 0

	key := strconv.FormatInt(id, 10)
	//err := conf.GetDbConn().Conn().Table(key + "_keys").Order(order).Find(&tss).Error
	err := conf.GetDbConn().Conn().Table("1_keys").Where("ecosystem = ?", key).Order(order).Find(&tss).Error
	if err != nil {
		return &ret, num, err
	}
	if page < 1 || size < 1 {
		return &ret, num, err
	}
	ioffet = (page - 1) * size
	num = int64(len(tss))
	if num < int64(page*size) {
		size = int(num) % size
	}
	if num < int64(ioffet) || num < 1 {
		return &ret, num, err
	}

	tokenSymbol, name := Tokens.Get(id), EcoNames.Get(id)
	if tokenSymbol == "" && name == "" {
		return &ret, num, err
	}
	for i = 0; i < size; i++ {
		da := KeyHex{}
		da.Ecosystem = id
		if da.Ecosystem == 1 {
			da.TokenSymbol = SysTokenSymbol
		} else {
			da.TokenSymbol = tokenSymbol
		}
		//da.Token_title = es.Token_title
		da.Ecosystemname = name
		da.ID = strconv.FormatInt(tss[ioffet].ID, 10)
		da.PublicKey = hex.EncodeToString(tss[ioffet].PublicKey)
		da.Maxpay = tss[ioffet].Maxpay
		da.Amount = tss[ioffet].Amount.String()
		da.Deleted = tss[ioffet].Deleted
		da.Multi = tss[ioffet].Multi
		da.Blocked = tss[ioffet].Blocked
		//fmt.Println("ecosystem %d", id)
		ret = append(ret, da)
		ioffet++
	}

	return &ret, num, err
}

func (m *Key) GetKeyDetail(keyId int64, wallet string) (*EcosyKeyTotalHex, error) {
	d := EcosyKeyTotalHex{}
	d.Ecosystem = m.Ecosystem
	d.Wallet = m.Account

	mb := Member{}
	fm, _ := mb.GetAccount(m.Ecosystem, wallet)
	if fm {
		if mb.ImageID != nil {
			if *mb.ImageID != int64(0) {
				hash, err := GetFileHash(*mb.ImageID)
				if err != nil {
					return &d, err
				}
				d.MemberHash = hash
			}
		}
	}
	if INameReady {
		ie := &IName{}
		f, err := ie.Get(wallet)
		if err == nil && f {
			d.MemberName = ie.Name
		}
	}

	//
	ems := Ecosystem{}
	f, err := ems.Get(m.Ecosystem)
	if err != nil {
		return &d, err
	}
	if f {
		d.Ecosystemname = ems.Name
		d.IsValued = ems.IsValued
		d.LogoHash = GetLogoHash(ems.ID)
		d.TokenSymbol = ems.TokenSymbol
		d.Digits = ems.Digits
	}

	ts := &History{}
	dh, err := ts.GetAccountHistoryTotals(m.Ecosystem, keyId)
	if err != nil {
		return &d, err
	}

	if AssignReady {
		ag := &AssignInfo{}
		lockAmount, err := ag.GetBalance(nil, wallet)
		if err != nil {
			return &d, err
		}
		d.LockAmount = d.LockAmount.Add(lockAmount)
	}
	if AirdropReady {
		airdrop := &AirdropInfo{}
		f, err = airdrop.Get(wallet)
		if err != nil {
			return nil, err
		}
		if f {
			d.LockAmount = d.LockAmount.Add(airdrop.BalanceAmount)
			d.StakeAmount = d.StakeAmount.Add(airdrop.StakeAmount)
		}
	}

	if m.Lock != "" {
		var stake KeyLocks
		if err := json.Unmarshal([]byte(m.Lock), &stake); err != nil {
			return &d, err
		}
		nftLock, _ := decimal.NewFromString(stake.NftMinerStake)
		referendumLock, _ := decimal.NewFromString(stake.CandidateReferendum)
		substituteLock, _ := decimal.NewFromString(stake.CandidateSubstitute)

		d.StakeAmount = d.StakeAmount.Add(nftLock.Add(referendumLock).Add(substituteLock))
	}

	d.Transaction = dh.Transaction
	d.InTx = dh.InTx
	d.OutTx = dh.OutTx
	d.InAmount = dh.InAmount
	d.OutAmount = dh.OutAmount
	d.TxAmount = d.InAmount.Add(d.OutAmount).String()

	accountAmount := decimal.Zero
	if m.Amount.GreaterThan(decimal.Zero) {
		accountAmount = m.Amount
	}
	var s1 SpentInfoHistory
	utxoAmount, err := s1.GetKeyBalance(keyId, 1)
	if err != nil {
		return nil, err
	}
	d.Amount = accountAmount.Add(utxoAmount).String()
	d.TotalAmount = d.TotalAmount.Add(accountAmount).Add(d.StakeAmount).Add(utxoAmount)
	d.JoinTime = getJoinTime(keyId, m.Ecosystem)
	d.RolesName = getRolesName(converter.AddressToString(keyId), m.Ecosystem)
	return &d, err
}

func getJoinTime(keyId int64, ecosystem int64) int64 {
	var tableId string
	ecosystemStr := strconv.FormatInt(ecosystem, 10)
	keyIdStr := strconv.FormatInt(keyId, 10)
	tableId = strings.Join([]string{keyIdStr, ecosystemStr}, ",")
	var rollback sqldb.RollbackTx
	f, err := isFound(GetDB(nil).Select("block_id").Table("rollback_tx").Where("table_name = '1_keys' AND table_id = ? AND data = ''", tableId).First(&rollback))
	if err != nil {
		log.Info("get join time err:", err.Error(), " table_id:", tableId)
		return 0
	}
	var bk Block
	if f {
		f, err := isFound(GetDB(nil).Select("time").Where("id = ?", rollback.BlockID).First(&bk))
		if err != nil {
			log.Info("get join time blockId err:", err.Error(), " blockId:", rollback.BlockID)
			return 0
		}

		if f {
			return bk.Time
		}
	} else {
		return FirstBlockTime
	}
	return 0
}

func getRolesName(account string, ecosystem int64) string {
	var roles sqldb.RolesParticipants
	var rolesName string
	var nameList []string
	roles.SetTablePrefix(ecosystem)
	list, err := roles.GetActiveMemberRoles(account)
	if err != nil {
		log.WithFields(log.Fields{"warn": err, "key": account}).Warn("GetActiveMemberRoles err")
		return ""
	}
	if len(list) == 0 {
		return ""
	}
	type roleData struct {
		Name string `json:"name"`
	}
	var role roleData
	for i := 0; i < len(list); i++ {
		if err := json.Unmarshal([]byte(list[i].Role), &role); err != nil {
			log.WithFields(log.Fields{"warn": err, "role": list[i].Role}).Warn("getRolesName json err")
			return ""
		}
		nameList = append(nameList, role.Name)
	}
	if len(nameList) > 0 {
		rolesName = strings.Join(nameList, " / ")
	}
	return rolesName
}

func (m *Key) GetWalletTotalBasisEcosystem(wallet string) (*EcosyKeyTotalHex, error) {
	var (
		ft  Key
		ret EcosyKeyTotalHex
	)
	wid := converter.StringToAddress(wallet)
	if wid != 0 || wallet == BlackHoleAddr {
		f, err := isFound(GetDB(nil).Table("1_keys").Where("id = ? and ecosystem = ?", wid, 1).First(&ft))
		if err != nil {
			return &ret, err
		}
		if !f {
			if wallet != BlackHoleAddr {
				var sp SpentInfo
				f, err = isFound(GetDB(nil).Where("output_key_id = ? AND ecosystem = ?", wid, 1).First(&sp))
				if err != nil {
					return nil, err
				}
				if !f {
					return nil, errors.New("account doesn't not exist")
				}
				ft.Ecosystem = 1
				ft.Account = wallet
			} else {
				ft.Ecosystem = 1
				ft.Account = BlackHoleAddr
			}
		}
		df, err := ft.GetKeyDetail(wid, wallet)
		if err != nil {
			return &ret, err
		}
		ret = *df
	} else {
		return &ret, errors.New("wallet invalid")
	}
	return &ret, nil
}

func GetWalletTokenChangeBasis(account string) (*AccountHistoryChart, error) {
	var (
		ret AccountHistoryChart
		his History
	)
	kid := converter.StringToAddress(account)
	if kid != 0 || account == BlackHoleAddr {
		chart, err := his.GetWalletTimeLineHistoryTotals(1, kid, 30)
		if err != nil {
			return &ret, err
		}
		return chart, nil
	} else {
		return &ret, errors.New("account invalid")
	}
}

// GetWalletTotalEcosystem response 10-20ms The modified 4-9ms
func (m *Key) GetWalletTotalEcosystem(page, limit int, order string, wallet string) (*GeneralResponse, error) {
	var (
		total int64
		ret   GeneralResponse
	)
	ret.Limit = limit
	ret.Page = page
	da := []EcosyKeyTotalDetail{}

	if page <= 0 || limit <= 0 {
		return nil, errors.New("request params invalid")
	}
	if order == "" {
		order = "total_amount desc"
	}

	wid := converter.StringToAddress(wallet)
	if wid != 0 || wallet == BlackHoleAddr {
		err := GetDB(nil).Model(AccountDetail{}).Where("id = ?", wid).Count(&total).Error
		if err != nil {
			return &ret, err
		}

		err = GetDB(nil).Table("account_detail AS ad").
			Select(`account,ecosystem,logo_hash,(SELECT array_to_string(array(SELECT rs.role->>'name' FROM "1_roles_participants" as rs 
	WHERE rs.ecosystem=ad.ecosystem and rs.member->>'account' = ad.account AND rs.deleted = 0),' / '))as roles_name,join_time,total_amount,stake_amount`).
			Where("id = ?", wid).Offset((page - 1) * limit).Limit(limit).Order(order).Find(&da).Error
		if err != nil {
			return &ret, err
		}

		ret.Total = int64(total)
		for k, v := range da {
			if v.Ecosystem == 1 {
				var as AssignInfo
				lockAmount, err := as.GetBalance(nil, v.Account)
				if err != nil {
					return nil, err
				}
				da[k].LockAmount = lockAmount

				if AirdropReady {
					airdrop := &AirdropInfo{}
					f, err := airdrop.Get(v.Account)
					if err != nil {
						return nil, err
					}
					if f {
						da[k].LockAmount = da[k].LockAmount.Add(airdrop.BalanceAmount)
						da[k].StakeAmount = da[k].StakeAmount.Add(airdrop.StakeAmount)
						da[k].TotalAmount = da[k].TotalAmount.Add(airdrop.StakeAmount)
					}
				}
			}
			da[k].Name = EcoNames.Get(v.Ecosystem)
			da[k].TokenSymbol = Tokens.Get(v.Ecosystem)
			da[k].Digits = EcoDigits.GetInt(v.Ecosystem, 0)
			da[k].Amount = da[k].TotalAmount.Sub(da[k].StakeAmount).String()
		}
		ret.List = da
		return &ret, nil
	} else {
		return &ret, errors.New("wallet invalid")
	}
}

func GetCirculations(ecosystem int64) (string, error) {
	var err error
	type result struct {
		Amount string
	}
	var (
		rets result
	)
	err = GetDB(nil).Raw(`
	SELECT coalesce(sum(amount),0)+
		(SELECT COALESCE(sum(output_value),0) FROM spent_info WHERE input_tx_hash is NULL AND ecosystem = ?) as amount 
	FROM "1_keys" WHERE ecosystem = ? AND id <> 0 AND id <> 5555
`, ecosystem, ecosystem).Take(&rets).Error
	if err != nil {
		return decimal.Zero.String(), err
	}

	return rets.Amount, err
}

func (m *Key) GetStakeAmount() (string, error) {
	type result struct {
		Amount decimal.Decimal
	}
	var agi AssignInfo
	agm, err := agi.GetBalance(nil, "")
	if err != nil {
		return "0", err
	}

	if HasTableOrView("1_mine_stake") {
		var mine, pool result
		err = conf.GetDbConn().Conn().Table("1_keys").Select("SUM(mine_lock) as amount").Where("ecosystem = 1").Scan(&mine).Error
		if err != nil {
			b := strings.ContainsAny(err.Error(), "column mine_lock does not exist")
			if b {
				return agm.String(), nil
			}
			return "0", err
		}
		err = conf.GetDbConn().Conn().Table("1_keys").Select("SUM(pool_lock) as amount").Where("ecosystem = 1").Scan(&pool).Error
		if err != nil {
			b := strings.ContainsAny(err.Error(), "column pool_lock does not exist")
			if b {
				return agm.String(), nil
			}
			return "0", err
		}

		rt := mine.Amount.Add(pool.Amount)
		rt = rt.Add(agm)
		return rt.String(), nil
	}
	return agm.String(), nil
}

func (m *Key) GetAccountList(page, limit int, ecosystem int64) (*KeysListResult, error) {

	var (
		ret KeysListResult
	)
	ret.Limit = limit
	ret.Page = page

	totalAmount, err := allKeyAmount.Get(ecosystem)
	if err != nil {
		return &ret, nil
	}

	err = GetDB(nil).Model(AccountDetail{}).Where("ecosystem = ?", ecosystem).Count(&ret.Total).Error
	if err != nil {
		return nil, err
	}
	if ecosystem == 1 && AirdropReady {
		err = GetDB(nil).Table("account_detail as ad").
			Select("ad.account,ad.total_amount+COALESCE(ai.stake_amount,0) AS total_amount,ad.stake_amount+COALESCE(ai.stake_amount,0) AS stake_amount").
			Joins(`LEFT JOIN "1_airdrop_info" AS ai ON(ai.account = ad.account)`).
			Where("ecosystem = ?", ecosystem).Offset((page - 1) * limit).Limit(limit).Order("total_amount desc,account ASC").Find(&ret.Rets).Error
	} else {
		err = GetDB(nil).Model(AccountDetail{}).
			Select("account,total_amount,stake_amount").
			Where("ecosystem = ?", ecosystem).Offset((page - 1) * limit).Limit(limit).Order("total_amount desc,account ASC").Find(&ret.Rets).Error
	}
	if err != nil {
		return nil, err
	}
	tokenSymbol := Tokens.Get(ecosystem)
	digits := EcoDigits.GetInt(ecosystem, 0)
	for i := 0; i < len(ret.Rets); i++ {
		ret.Rets[i].Ecosystem = ecosystem
		if ecosystem == 1 {
			if AssignReady {
				ag := &AssignInfo{}
				ba, err := ag.GetBalance(nil, ret.Rets[i].Account)
				if err != nil {
					return &ret, err
				}
				ret.Rets[i].LockAmount = ret.Rets[i].LockAmount.Add(ba)
			}
			if AirdropReady {
				air := &AirdropInfo{}
				f, err := air.Get(ret.Rets[i].Account)
				if err == nil && f {
					ret.Rets[i].LockAmount = ret.Rets[i].LockAmount.Add(air.BalanceAmount)
				}
			}
		}
		ret.Rets[i].TokenSymbol = tokenSymbol
		ret.Rets[i].Digits = digits
		ret.Rets[i].Amount = ret.Rets[i].TotalAmount.Sub(ret.Rets[i].StakeAmount).String()

		amount, _ := decimal.NewFromString(ret.Rets[i].Amount)
		if amount.GreaterThan(decimal.Zero) {
			ret.Rets[i].AccountedFor = amount.Mul(decimal.NewFromInt(100)).DivRound(totalAmount, 2)
		}
		if INameReady {
			in := &IName{}
			f, err := in.Get(ret.Rets[i].Account)
			if err == nil && f {
				ret.Rets[i].AccountName = in.Name
			}
		}
	}

	return &ret, nil
}

func GetAccountListChart(ecosystem int64) (*KeysListChartResult, error) {
	var (
		ret KeysListChartResult
	)

	if ecosystem == 1 {
		ret.KeysInfo, _ = getScanOutKeyInfoFromRedis()
	} else {
		ret.KeysInfo, _ = getScanOutKeyInfo(ecosystem)
	}
	newKey, err := Get15DaysNewKeyFromRedis(ecosystem)
	if err != nil {
		return &ret, err
	}
	ret.KeyChart = *newKey

	return &ret, nil
}

func (k *Key) GetEcosystemTokenSymbolList(page, limit int, ecosystem int64) (*GeneralResponse, error) {
	var (
		rets GeneralResponse
		err  error
	)
	rets.Limit = limit
	rets.Page = page

	ecoTotal, err := allKeyAmount.Get(ecosystem)
	if err != nil {
		return &rets, nil
	}
	tokenSymbol := Tokens.Get(ecosystem)
	digits := EcoDigits.GetInt(ecosystem, 0)

	var ret []EcosystemTokenSymbolList
	var querySql *gorm.DB
	if ecosystem == 1 {
		if AirdropReady {
			querySql = GetDB(nil).Table(`"account_detail" as ad`).Select(`id,account,ecosystem,total_amount +
			COALESCE((SELECT stake_amount FROM "1_airdrop_info" WHERE account = ad.account),0) as amount`).Where("ecosystem = 1 AND total_amount > 0").Offset((page - 1) * limit).Limit(limit)
		} else {
			querySql = GetDB(nil).Model(AccountDetail{}).Select(`id,account,ecosystem,total_amount as amount`).Where("ecosystem = 1 AND total_amount > 0").Offset((page - 1) * limit).Limit(limit)
		}
	} else {
		querySql = GetDB(nil).Raw(`
SELECT v1.id,CASE WHEN v1.join_time > 0 THEN TRUE ELSE FALSE END AS activation,v1.account,v1.ecosystem,v1.total_amount AS amount,
	CASE WHEN (SELECT control_mode FROM "1_ecosystems" WHERE id = v1.ecosystem AND id > 1) = 2 THEN
		CASE WHEN (
			SELECT count(1) FROM "1_votings_participants" WHERE
				voting_id = (SELECT id FROM "1_votings" WHERE deleted = 0 AND voting->>'name' like '%voting_for_control_mode_template%' AND ecosystem = v1.ecosystem ORDER BY id DESC LIMIT 1)
							AND member->>'account'=v1.account) > 0 
		THEN
			TRUE
		ELSE
			FALSE
		END
	ELSE
	 FALSE
	END AS front_committee
FROM(
	SELECT id,account,ecosystem,total_amount,join_time FROM "account_detail" WHERE ecosystem = ? AND total_amount > 0
)AS v1
ORDER BY amount desc OFFSET ? LIMIT ?
`, ecosystem, (page-1)*limit, limit)
	}
	err = GetDB(nil).Model(AccountDetail{}).Where("ecosystem = ? AND total_amount > 0", ecosystem).Count(&rets.Total).Error
	if err != nil {
		return nil, err
	}
	err = querySql.Find(&ret).Error
	if err != nil {
		return nil, err
	}
	type ids struct {
		Id int64
	}
	var committeeList []ids
	if ecosystem > 1 {
		err = GetDB(nil).Raw(`
select
  id
from
  (
  select sum(amount) as amount,
    id,CASE 
    WHEN true THEN dense_rank() OVER (ORDER BY sum ( amount )  DESC)
    ELSE rank() OVER (ORDER BY sum (amount)  DESC) END AS rank
  
  
   from
    (
    select
      output_value as amount,
      output_key_id as id 
    from
      spent_info
      left join "1_keys" on spent_info.ecosystem = "1_keys".ecosystem 
      and id = spent_info.output_key_id 
    where
      input_tx_hash is null 
      and "spent_info".ecosystem = ?
      and blocked = 0 
      and deleted = 0 
      and length ( pub ) > 0 union
    select
      amount,
    id 
    from
      "1_keys" 
    where
      ecosystem = ? 
      and blocked = 0 
      and deleted = 0 
      and amount > 0 
      and length ( pub ) > 0 
    ) tmp 
  group by
  id 
  order by
  rank asc 
  ) r 
where
  rank <= 50
order by
  rank,
  id;
`, ecosystem, ecosystem).Find(&committeeList).Error
		if err != nil {
			return nil, err
		}
	}
	for i := 0; i < len(ret); i++ {
		for _, v := range committeeList {
			if v.Id == ret[i].Id {
				ret[i].Committee = true
			}
		}
		if INameReady {
			ie := &IName{}
			f, err := ie.Get(ret[i].Account)
			if err == nil && f {
				ret[i].AccountName = ie.Name
			}
		}
		ret[i].TokenSymbol = tokenSymbol
		ret[i].Digits = digits
		ret[i].AccountedFor = ret[i].Amount.Mul(decimal.NewFromInt(100)).DivRound(ecoTotal, 2)
	}
	rets.List = ret

	return &rets, err
}

func getEcosystemNewKeyChart(ecosystem int64, getDay int) KeyInfoChart {
	var rets KeyInfoChart
	tz := time.Unix(GetNowTimeUnix(), 0)
	yesterday := time.Date(tz.Year(), tz.Month(), tz.Day()-1, 0, 0, 0, 0, tz.Location())
	t1 := yesterday.AddDate(0, 0, -1*getDay)

	rets.Time = make([]int64, getDay)
	rets.NewKey = make([]int64, getDay)

	var keyList []DaysNumber
	err := GetDB(nil).Model(AccountDetail{}).Select("to_char(to_timestamp(join_time),'yyyy-MM-dd') days,count(1) num").Group("days").
		Where("ecosystem = ? AND join_time >= ?", ecosystem, t1.Unix()).Find(&keyList).Error
	if err != nil {
		log.WithFields(log.Fields{"error": err, "t1": t1.Unix()}).Warn("get Ecosystem New Key Chart Failed")
		return rets
	}

	for i := 0; i < len(rets.Time); i++ {
		rets.Time[i] = t1.AddDate(0, 0, i+1).Unix()
		rets.NewKey[i] = GetDaysNumber(rets.Time[i], keyList)
	}
	rets.Name = EcoNames.Get(ecosystem)

	return rets
}

func GetAccountTotalAmount(ecosystem int64, account string) (AccountTotalAmountChart, error) {
	var (
		rets AccountTotalAmountChart
		err  error
		f    bool
	)
	keyId := converter.StringToAddress(account)
	if keyId == 0 {
		return rets, errors.New("account invalid:" + account + " in ecosystem:" + strconv.FormatInt(ecosystem, 10))
	}
	if (NftMinerReady || NodeReady) && ecosystem == 1 {
		f, err = isFound(GetDB(nil).Raw(`
			SELECT sum(amount) AS amount,sum(stake_amount) AS stake_amount FROM(
				SELECT k1.amount,
				to_number(coalesce(NULLIF(k1.lock->>'nft_miner_stake',''),'0'),'999999999999999999999999999999')+ 
						to_number(coalesce(NULLIF(k1.lock->>'candidate_referendum',''),'0'),'999999999999999999999999999999') + 
						to_number(coalesce(NULLIF(k1.lock->>'candidate_substitute',''),'0'),'999999999999999999999999999999')as stake_amount FROM "1_keys" AS k1 WHERE ecosystem = ? and id = ?
				UNION 
				SELECT sum(output_value)AS amount,0 AS stake_amount FROM spent_info WHERE input_tx_hash is NULL AND ecosystem = ? AND output_key_id = ?
			)AS v1
`, ecosystem, keyId, ecosystem, keyId).Take(&rets))

	} else {
		f, err = isFound(GetDB(nil).Raw(`
			SELECT sum(amount) AS amount,sum(stake_amount) AS stake_amount FROM(
				SELECT k1.amount,
				0 AS stake_amount FROM "1_keys" AS k1 WHERE ecosystem = ? and id = ?
				UNION 
				SELECT sum(output_value)AS amount,0 AS stake_amount FROM spent_info WHERE input_tx_hash is NULL AND ecosystem = ? AND output_key_id = ?		
			)AS v1
`, ecosystem, keyId, ecosystem, keyId).Take(&rets))
	}
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Account Total Amount Chart Failed")
		return rets, nil
	}
	if !f {
		return rets, errors.New("unknown account:" + account + " in ecosystem:" + strconv.FormatInt(ecosystem, 10))
	}
	if AssignReady {
		ais := AssignInfo{}
		lockAmount, err := ais.GetBalance(nil, account)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Get Account Total Amount Chart lock amount Failed")
			return rets, nil
		}
		rets.LockAmount = lockAmount

	}
	if AirdropReady {
		airdrop := &AirdropInfo{}
		f, err = airdrop.Get(account)
		if err != nil {
			return rets, err
		}
		if f {
			rets.LockAmount = rets.LockAmount.Add(airdrop.BalanceAmount)
			rets.StakeAmount = rets.StakeAmount.Add(airdrop.StakeAmount)
		}
	}

	rets.TokenSymbol = Tokens.Get(ecosystem)
	rets.Digits = EcoDigits.GetInt(ecosystem, 0)
	rets.TotalAmount = rets.Amount.Add(rets.LockAmount).Add(rets.StakeAmount)
	zeroDec := decimal.New(0, 0)
	if rets.TotalAmount.GreaterThan(zeroDec) {
		if rets.Amount.GreaterThan(zeroDec) {
			rets.AmountRatio, _ = rets.Amount.Mul(decimal.NewFromInt(100)).DivRound(rets.TotalAmount, 2).Float64()
		}
		if rets.StakeAmount.GreaterThan(zeroDec) {
			rets.StakeAmountRatio, _ = rets.StakeAmount.Mul(decimal.NewFromInt(100)).DivRound(rets.TotalAmount, 2).Float64()
		}
		if rets.LockAmount.GreaterThan(zeroDec) {
			rets.LockRatio, _ = rets.LockAmount.Mul(decimal.NewFromInt(100)).DivRound(rets.TotalAmount, 2).Float64()
		}
	}
	return rets, nil
}
