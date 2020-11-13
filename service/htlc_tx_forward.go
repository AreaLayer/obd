package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/omnilaboratory/obd/bean"
	"github.com/omnilaboratory/obd/bean/enum"
	"github.com/omnilaboratory/obd/config"
	"github.com/omnilaboratory/obd/dao"
	"github.com/omnilaboratory/obd/rpc"
	"github.com/omnilaboratory/obd/tool"
	trackerBean "github.com/omnilaboratory/obd/tracker/bean"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type htlcForwardTxManager struct {
	operationFlag sync.Mutex

	//缓存数据
	addHtlcTempDataAtBobSide map[string]string
	htlcInvoiceTempData      map[string]bean.HtlcRequestFindPathInfo

	//在步骤1，缓存需要发往40号协议的信息
	tempDataSendTo40PAtAliceSide map[string]bean.CreateHtlcTxForC3aOfP2p
	//在步骤4，缓存需要发往41号协议的信息
	tempDataSendTo41PAtBobSide map[string]bean.NeedAliceSignHtlcTxOfC3bP2p
	//在步骤6，缓存来自41号协议的信息
	tempDataFrom41PAtAliceSide map[string]bean.NeedAliceSignHtlcTxOfC3bP2p
	//在步骤7，缓存需要发往42号协议的信息
	tempDataSendTo42PAtAliceSide map[string]bean.NeedBobSignHtlcSubTxOfC3bP2p
	//在步骤9，缓存来自42号协议的信息
	tempDataFrom42PAtBobSide map[string]bean.NeedBobSignHtlcSubTxOfC3bP2p
}

// htlc pay money  付款
var HtlcForwardTxService htlcForwardTxManager

func (service *htlcForwardTxManager) CreateHtlcInvoice(msg bean.RequestMessage, user bean.User) (data interface{}, err error) {

	if tool.CheckIsString(&msg.Data) == false {
		return nil, errors.New(enum.Tips_common_empty + "msd data")
	}

	requestData := &bean.HtlcRequestInvoice{}
	err = json.Unmarshal([]byte(msg.Data), requestData)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	addr := ""
	//obbc,obtb,obcrt
	if strings.Contains(config.ChainNode_Type, "main") {
		addr = "obbc"
	}
	if strings.Contains(config.ChainNode_Type, "test") {
		addr = "obtb"
	}
	if strings.Contains(config.ChainNode_Type, "reg") {
		addr = "obcrt"
	}
	if requestData.Amount < config.GetOmniDustBtc() {
		return nil, errors.New(enum.Tips_common_wrong + "amount")
	} else {
		requestData.Amount *= 100000000
		temp := int(requestData.Amount)
		addr += strconv.Itoa(temp) + "s"
	}

	addr += "1"

	if requestData.PropertyId < 0 {
		return nil, errors.New(enum.Tips_common_wrong + "property_id")
	}
	_, err = rpcClient.OmniGetProperty(requestData.PropertyId)
	if err != nil {
		return nil, err
	} else {
		propertyId := ""
		tool.ConvertNumToString(int(requestData.PropertyId), &propertyId)
		code, err := tool.GetMsgLengthFromInt(len(propertyId))
		if err != nil {
			return nil, err
		}
		addr += "p" + code + propertyId
	}

	code, err := tool.GetMsgLengthFromInt(len(msg.SenderNodePeerId))
	if err != nil {
		return nil, err
	}
	addr += "n" + code + msg.SenderNodePeerId

	code, err = tool.GetMsgLengthFromInt(len(msg.SenderUserPeerId))
	if err != nil {
		return nil, err
	}
	addr += "u" + code + msg.SenderUserPeerId

	if tool.CheckIsString(&requestData.H) == false {
		return nil, errors.New(enum.Tips_common_wrong + "h")
	} else {
		//ph payment H
		code, err = tool.GetMsgLengthFromInt(len(requestData.H))
		if err != nil {
			return nil, err
		}
		addr += "h" + code + requestData.H
	}

	if time.Time(requestData.ExpiryTime).IsZero() {
		return nil, errors.New(enum.Tips_common_wrong + "expiry_time")
	} else {
		if time.Now().After(time.Time(requestData.ExpiryTime)) {
			return nil, errors.New(fmt.Sprintf(enum.Tips_htlc_expiryTimeAfterNow, "expiry_time"))
		}
		expiryTime := ""
		tool.ConvertNumToString(int(time.Time(requestData.ExpiryTime).Unix()), &expiryTime)
		code, err = tool.GetMsgLengthFromInt(len(expiryTime))
		if err != nil {
			return nil, err
		}
		addr += "x" + code + expiryTime
	}

	code, err = tool.GetMsgLengthFromInt(1)
	if err != nil {
		return nil, err
	}
	isPrivate := "0"
	if requestData.IsPrivate {
		isPrivate = "1"
	}
	addr += "t" + code + isPrivate

	if len(requestData.Description) > 0 {
		code, err = tool.GetMsgLengthFromInt(len(requestData.Description))
		if err != nil {
			return nil, err
		}
		addr += "d" + code + requestData.Description
	}

	bytes := []byte(addr)
	sum := 0
	for _, item := range bytes {
		sum += int(item)
	}
	checkSum := ""
	tool.ConvertNumToString(sum, &checkSum)

	addr += checkSum

	return addr, nil
}

// 401 find htlc find path
func (service *htlcForwardTxManager) PayerRequestFindPath(msgData string, user bean.User) (data interface{}, isPrivate bool, err error) {
	if tool.CheckIsString(&msgData) == false {
		return nil, false, errors.New(enum.Tips_common_empty + "msg data")
	}

	requestData := &bean.HtlcRequestFindPath{}
	err = json.Unmarshal([]byte(msgData), requestData)
	if err != nil {
		log.Println(err.Error())
		return nil, false, err
	}

	var requestFindPathInfo bean.HtlcRequestFindPathInfo

	if tool.CheckIsString(&requestData.Invoice) {
		htlcRequestInvoice, err := tool.DecodeInvoiceObjFromCodes(requestData.Invoice)
		if err != nil {
			return nil, false, errors.New(enum.Tips_common_wrong + "invoice")
		}
		if err = findUserIsOnline(htlcRequestInvoice.RecipientNodePeerId, htlcRequestInvoice.RecipientUserPeerId); err != nil {
			return nil, requestFindPathInfo.IsPrivate, err
		}
		requestFindPathInfo = htlcRequestInvoice.HtlcRequestFindPathInfo
	} else {
		requestFindPathInfo = requestData.HtlcRequestFindPathInfo
		if tool.CheckIsString(&requestFindPathInfo.RecipientNodePeerId) == false {
			return nil, requestFindPathInfo.IsPrivate, errors.New(enum.Tips_common_wrong + "recipient_node_peer_id")
		}
		if tool.CheckIsString(&requestFindPathInfo.RecipientUserPeerId) == false {
			return nil, requestFindPathInfo.IsPrivate, errors.New(enum.Tips_common_wrong + "recipient_user_peer_id")
		}

		if err = findUserIsOnline(requestFindPathInfo.RecipientNodePeerId, requestFindPathInfo.RecipientUserPeerId); err != nil {
			return nil, requestFindPathInfo.IsPrivate, err
		}
	}

	if requestFindPathInfo.PropertyId < 0 {
		return nil, requestFindPathInfo.IsPrivate, errors.New(enum.Tips_common_wrong + "property_id")
	}

	_, err = rpcClient.OmniGetProperty(requestFindPathInfo.PropertyId)
	if err != nil {
		return nil, requestFindPathInfo.IsPrivate, err
	}

	if requestFindPathInfo.Amount < config.GetOmniDustBtc() {
		return nil, requestFindPathInfo.IsPrivate, errors.New(enum.Tips_common_wrong + "amount")
	}

	if time.Now().After(time.Time(requestFindPathInfo.ExpiryTime)) {
		return nil, requestFindPathInfo.IsPrivate, errors.New(fmt.Sprintf(enum.Tips_htlc_expiryTimeAfterNow, "expiry_time"))
	}

	if requestFindPathInfo.IsPrivate == false {
		//tracker find path
		pathRequest := trackerBean.HtlcPathRequest{}
		pathRequest.H = requestFindPathInfo.H
		pathRequest.PropertyId = requestFindPathInfo.PropertyId
		pathRequest.Amount = requestFindPathInfo.Amount
		pathRequest.RealPayerPeerId = user.PeerId
		pathRequest.PayerObdNodeId = tool.GetObdNodeId()
		pathRequest.PayeePeerId = requestFindPathInfo.RecipientUserPeerId
		sendMsgToTracker(enum.MsgType_Tracker_GetHtlcPath_351, pathRequest)
		if service.htlcInvoiceTempData == nil {
			service.htlcInvoiceTempData = make(map[string]bean.HtlcRequestFindPathInfo)
		}
		service.htlcInvoiceTempData[user.PeerId+"_"+pathRequest.H] = requestFindPathInfo
		return make(map[string]interface{}), requestFindPathInfo.IsPrivate, nil
	} else {
		requestData.HtlcRequestFindPathInfo = requestFindPathInfo
		return getPrivateChannelForHtlc(requestData, user)
	}
}

func getPrivateChannelForHtlc(requestData *bean.HtlcRequestFindPath, user bean.User) (data interface{}, isPrivate bool, err error) {
	tx, err := user.Db.Begin(true)
	if err != nil {
		log.Println(err)
		return nil, true, err
	}
	defer tx.Rollback()
	//get all private channel
	var nodes []dao.ChannelInfo
	err = tx.Select(
		q.Eq("PropertyId", requestData.PropertyId),
		q.Eq("IsPrivate", true),
		q.Eq("CurrState", dao.ChannelState_CanUse),
		q.Or(
			q.And(
				q.Eq("PeerIdB", requestData.RecipientUserPeerId),
				q.Eq("PeerIdA", user.PeerId)),
			q.And(
				q.Eq("PeerIdB", user.PeerId),
				q.Eq("PeerIdA", requestData.RecipientUserPeerId)),
		)).OrderBy("CreateAt").Reverse().Find(&nodes)

	retData := make(map[string]interface{})
	if nodes != nil && len(nodes) > 0 {
		for _, channel := range nodes {
			commitmentTxInfo, err := getLatestCommitmentTxUseDbTx(tx, channel.ChannelId, user.PeerId)
			if err == nil && commitmentTxInfo.Id > 0 {
				if commitmentTxInfo.AmountToRSMC >= requestData.Amount {
					retData["h"] = requestData.H
					retData["is_private"] = requestData.IsPrivate
					retData["property_id"] = requestData.PropertyId
					retData["amount"] = requestData.Amount
					retData["routing_packet"] = channel.ChannelId
					retData["min_cltv_expiry"] = 1
					retData["next_node_peerId"] = requestData.RecipientUserPeerId
					retData["memo"] = requestData.Description
					break
				}
			}
		}
	}
	_ = tx.Commit()
	if len(retData) == 0 {
		return nil, true, errors.New(enum.Tips_htlc_noPrivatePath)
	}
	return retData, true, nil
}

func (service *htlcForwardTxManager) GetResponseFromTrackerOfPayerRequestFindPath(channelPath string, user bean.User) (data interface{}, err error) {
	if tool.CheckIsString(&channelPath) == false {
		err = errors.New("has no channel path")
		log.Println(err)
		return nil, err
	}

	log.Println(channelPath)

	dataArr := strings.Split(channelPath, "_")
	if len(dataArr) != 3 {
		return nil, errors.New("no channel path")
	}

	h := dataArr[0]
	requestFindPathInfo := service.htlcInvoiceTempData[user.PeerId+"_"+h]
	if &requestFindPathInfo == nil {
		return nil, errors.New("has no channel path")
	}

	splitArr := strings.Split(dataArr[1], ",")
	currChannelInfo := dao.ChannelInfo{}
	err = user.Db.Select(
		q.Eq("ChannelId", splitArr[0]),
		q.Eq("CurrState", dao.ChannelState_CanUse),
		q.Or(
			q.Eq("PeerIdA", user.PeerId),
			q.Eq("PeerIdB", user.PeerId))).First(&currChannelInfo)
	if err != nil {
		err = errors.New("has no ChannelPath")
		log.Println(err)
		return nil, err
	}
	nextNodePeerId := currChannelInfo.PeerIdB
	if user.PeerId == currChannelInfo.PeerIdB {
		nextNodePeerId = currChannelInfo.PeerIdA
	}

	arrLength := len(strings.Split(dataArr[1], ","))
	retData := make(map[string]interface{})
	retData["h"] = h
	retData["is_private"] = false
	retData["property_id"] = requestFindPathInfo.PropertyId
	retData["amount"] = requestFindPathInfo.Amount
	retData["routing_packet"] = dataArr[1]
	retData["min_cltv_expiry"] = arrLength
	retData["next_node_peerId"] = nextNodePeerId
	retData["memo"] = requestFindPathInfo.Description

	delete(service.htlcInvoiceTempData, user.PeerId+"_"+h)
	return retData, nil
}

// step 1 alice -100040协议的alice方的逻辑 alice start a htlc as payer
func (service *htlcForwardTxManager) AliceAddHtlcAtAliceSide(msg bean.RequestMessage, user bean.User) (data interface{}, needSign bool, err error) {
	if tool.CheckIsString(&msg.Data) == false {
		return nil, false, errors.New(enum.Tips_common_empty + "msg data")
	}

	requestData := &bean.CreateHtlcTxForC3a{}
	err = json.Unmarshal([]byte(msg.Data), requestData)
	if err != nil {
		log.Println(err.Error())
		return nil, false, err
	}

	tx, err := user.Db.Begin(true)
	if err != nil {
		log.Println(err)
		return nil, false, err
	}
	defer tx.Rollback()

	//region check input data 检测输入输入数据
	if requestData.Amount < config.GetOmniDustBtc() {
		return nil, false, errors.New(fmt.Sprintf(enum.Tips_common_amountMustGreater, config.GetOmniDustBtc()))
	}
	if tool.CheckIsString(&requestData.H) == false {
		return nil, false, errors.New(enum.Tips_common_empty + "h")
	}
	if tool.CheckIsString(&requestData.RoutingPacket) == false {
		return nil, false, errors.New(enum.Tips_common_empty + "routing_packet")
	}

	channelIds := strings.Split(requestData.RoutingPacket, ",")
	totalStep := len(channelIds)
	var channelInfo *dao.ChannelInfo
	var currStep = 0
	for index, channelId := range channelIds {
		temp := getChannelInfoByChannelId(tx, channelId, user.PeerId)
		if temp != nil {
			if temp.PeerIdA == msg.RecipientUserPeerId || temp.PeerIdB == msg.RecipientUserPeerId {
				channelInfo = temp
				currStep = index
				break
			}
		}
	}
	if channelInfo == nil {
		return nil, false, errors.New(enum.Tips_htlc_noChanneFromRountingPacket)
	}

	if channelInfo.CurrState == dao.ChannelState_NewTx {
		return nil, false, errors.New(enum.Tips_common_newTxMsg)
	}

	fundingTransaction := getFundingTransactionByChannelId(tx, channelInfo.ChannelId, user.PeerId)
	duration := time.Now().Sub(fundingTransaction.CreateAt)
	if duration > time.Minute*30 {
		if checkChannelOmniAssetAmount(*channelInfo) == false {
			err = errors.New(enum.Tips_rsmc_broadcastedChannel)
			log.Println(err)
			return nil, false, err
		}
	}

	if requestData.CltvExpiry < (totalStep - currStep) {
		requestData.CltvExpiry = totalStep - currStep
	}

	err = checkBtcFundFinish(channelInfo.ChannelAddress, false)
	if err != nil {
		log.Println(err)
		return nil, false, err
	}

	if tool.CheckIsString(&requestData.LastTempAddressPrivateKey) == false {
		err = errors.New(enum.Tips_common_empty + "last_temp_address_private_key")
		log.Println(err)
		return nil, false, err
	}

	latestCommitmentTx, _ := getLatestCommitmentTxUseDbTx(tx, channelInfo.ChannelId, user.PeerId)
	if latestCommitmentTx.Id > 0 && latestCommitmentTx.CurrState == dao.TxInfoState_Init {
		tx.DeleteStruct(latestCommitmentTx)
	}
	latestCommitmentTx, _ = getLatestCommitmentTxUseDbTx(tx, channelInfo.ChannelId, user.PeerId)

	if latestCommitmentTx.Id > 0 {
		if latestCommitmentTx.CurrState == dao.TxInfoState_CreateAndSign {
			_, err = tool.GetPubKeyFromWifAndCheck(requestData.LastTempAddressPrivateKey, latestCommitmentTx.RSMCTempAddressPubKey)
			if err != nil {
				return nil, false, errors.New(fmt.Sprintf(enum.Tips_rsmc_wrongPrivateKeyForLast, requestData.LastTempAddressPrivateKey, latestCommitmentTx.RSMCTempAddressPubKey))
			}
		}
		if latestCommitmentTx.CurrState == dao.TxInfoState_Create {
			if latestCommitmentTx.TxType != dao.CommitmentTransactionType_Htlc {
				return nil, false, errors.New("error commitment tx type")
			}

			if requestData.CurrRsmcTempAddressPubKey != latestCommitmentTx.RSMCTempAddressPubKey {
				return nil, false, errors.New(fmt.Sprintf(enum.Tips_rsmc_notSameValueWhenCreate, requestData.CurrRsmcTempAddressPubKey, latestCommitmentTx.RSMCTempAddressPubKey))
			}

			if requestData.CurrHtlcTempAddressPubKey != latestCommitmentTx.HTLCTempAddressPubKey {
				return nil, false, errors.New(fmt.Sprintf(enum.Tips_rsmc_notSameValueWhenCreate, requestData.CurrHtlcTempAddressPubKey, latestCommitmentTx.HTLCTempAddressPubKey))
			}

			if latestCommitmentTx.LastCommitmentTxId > 0 {
				lastCommitmentTx := &dao.CommitmentTransaction{}
				_ = tx.One("Id", latestCommitmentTx.LastCommitmentTxId, lastCommitmentTx)
				_, err = tool.GetPubKeyFromWifAndCheck(requestData.LastTempAddressPrivateKey, lastCommitmentTx.RSMCTempAddressPubKey)
				if err != nil {
					return nil, false, errors.New(fmt.Sprintf(enum.Tips_rsmc_wrongPrivateKeyForLast, requestData.LastTempAddressPrivateKey, lastCommitmentTx.RSMCTempAddressPubKey))
				}
			}
		}
	}
	if tool.CheckIsString(&requestData.CurrRsmcTempAddressPubKey) == false {
		err = errors.New(enum.Tips_common_empty + "curr_rsmc_temp_address_pub_key")
		log.Println(err)
		return nil, false, err
	}

	if tool.CheckIsString(&requestData.CurrHtlcTempAddressPubKey) == false {
		err = errors.New(enum.Tips_common_empty + "curr_htlc_temp_address_pub_key")
		log.Println(err)
		return nil, false, err
	}

	if tool.CheckIsString(&requestData.CurrHtlcTempAddressForHt1aPubKey) == false {
		err = errors.New(enum.Tips_common_empty + "curr_htlc_temp_address_for_ht1a_pub_key")
		log.Println(err)
		return nil, false, err
	}
	//endregion

	//这次请求的第一次发起
	htlcRequestInfo := &dao.AddHtlcRequestInfo{}
	_ = tx.Select(
		q.Eq("ChannelId", channelInfo.ChannelId),
		q.Eq("PropertyId", channelInfo.PropertyId),
		q.Eq("H", requestData.H),
		q.Eq("Amount", requestData.Amount),
		q.Eq("RoutingPacket", requestData.RoutingPacket),
		q.Eq("RecipientUserPeerId", msg.RecipientUserPeerId)).First(htlcRequestInfo)

	if htlcRequestInfo.Id == 0 || latestCommitmentTx.Id == 0 || latestCommitmentTx.CurrState == dao.TxInfoState_CreateAndSign {
		htlcRequestInfo.RecipientUserPeerId = msg.RecipientUserPeerId
		htlcRequestInfo.H = requestData.H
		htlcRequestInfo.Memo = requestData.Memo
		htlcRequestInfo.PropertyId = channelInfo.PropertyId
		htlcRequestInfo.Amount = requestData.Amount
		htlcRequestInfo.ChannelId = channelInfo.ChannelId
		htlcRequestInfo.RoutingPacket = requestData.RoutingPacket
		htlcRequestInfo.CurrRsmcTempAddressPubKey = requestData.CurrRsmcTempAddressPubKey
		htlcRequestInfo.CurrHtlcTempAddressPubKey = requestData.CurrHtlcTempAddressPubKey
		htlcRequestInfo.CurrHtlcTempAddressForHt1aIndex = requestData.CurrHtlcTempAddressForHt1aIndex
		htlcRequestInfo.CurrHtlcTempAddressForHt1aPubKey = requestData.CurrHtlcTempAddressForHt1aPubKey
		htlcRequestInfo.CurrState = dao.NS_Create
		htlcRequestInfo.CreateAt = time.Now()
		htlcRequestInfo.CreateBy = user.PeerId
		_ = tx.Save(htlcRequestInfo)

		totalStep := len(channelIds)
		latestCommitmentTx, err = htlcPayerCreateCommitmentTx_C3a(tx, channelInfo, *requestData, totalStep, currStep, latestCommitmentTx, user)
		if err != nil {
			log.Println(err)
			return nil, false, err
		}

		//更新tracker的htlc的状态
		if channelInfo.IsPrivate == false {
			txStateRequest := trackerBean.UpdateHtlcTxStateRequest{}
			txStateRequest.Path = latestCommitmentTx.HtlcRoutingPacket
			txStateRequest.H = latestCommitmentTx.HtlcH
			txStateRequest.DirectionFlag = trackerBean.HtlcTxState_PayMoney
			txStateRequest.CurrChannelId = channelInfo.ChannelId
			sendMsgToTracker(enum.MsgType_Tracker_UpdateHtlcTxState_352, txStateRequest)
		}
	} else {
		if requestData.CurrHtlcTempAddressForHt1aPubKey != htlcRequestInfo.CurrHtlcTempAddressForHt1aPubKey {
			return nil, false, errors.New(fmt.Sprintf(enum.Tips_rsmc_notSameValueWhenCreate, requestData.CurrHtlcTempAddressForHt1aPubKey, htlcRequestInfo.CurrHtlcTempAddressForHt1aPubKey))
		}
	}
	_ = tx.Commit()

	c3aP2pData := &bean.CreateHtlcTxForC3aOfP2p{}
	c3aP2pData.RoutingPacket = requestData.RoutingPacket
	c3aP2pData.ChannelId = channelInfo.ChannelId
	c3aP2pData.H = requestData.H
	c3aP2pData.Amount = requestData.Amount
	c3aP2pData.Memo = requestData.Memo
	c3aP2pData.CltvExpiry = requestData.CltvExpiry
	c3aP2pData.LastTempAddressPrivateKey = requestData.LastTempAddressPrivateKey
	c3aP2pData.CurrRsmcTempAddressPubKey = requestData.CurrRsmcTempAddressPubKey
	c3aP2pData.CurrHtlcTempAddressPubKey = requestData.CurrHtlcTempAddressPubKey
	c3aP2pData.C3aRsmcPartialSignedData = latestCommitmentTx.RsmcRawTxData
	c3aP2pData.C3aHtlcPartialSignedData = latestCommitmentTx.HtlcRawTxData
	c3aP2pData.C3aCounterpartyPartialSignedData = latestCommitmentTx.ToCounterpartyRawTxData
	c3aP2pData.CurrHtlcTempAddressForHt1aPubKey = requestData.CurrHtlcTempAddressForHt1aPubKey
	c3aP2pData.PayerCommitmentTxHash = latestCommitmentTx.CurrHash
	c3aP2pData.PayerNodeAddress = msg.SenderNodePeerId
	c3aP2pData.PayerPeerId = msg.SenderUserPeerId

	if latestCommitmentTx.CurrState == dao.TxInfoState_Init {
		txForC3a := bean.NeedAliceSignCreateHtlcTxForC3a{}
		txForC3a.ChannelId = latestCommitmentTx.ChannelId
		txForC3a.C3aRsmcRawData = latestCommitmentTx.RsmcRawTxData
		txForC3a.C3aHtlcRawData = latestCommitmentTx.HtlcRawTxData
		txForC3a.C3aCounterpartyRawData = latestCommitmentTx.ToCounterpartyRawTxData
		txForC3a.PayerNodeAddress = msg.SenderNodePeerId
		txForC3a.PayerPeerId = msg.SenderUserPeerId

		if service.tempDataSendTo40PAtAliceSide == nil {
			service.tempDataSendTo40PAtAliceSide = make(map[string]bean.CreateHtlcTxForC3aOfP2p)
		}
		service.tempDataSendTo40PAtAliceSide[user.PeerId+"_"+channelInfo.ChannelId] = *c3aP2pData
		return txForC3a, true, nil
	}
	return c3aP2pData, false, nil
}

// step 2 alice -100100 Alice对C3a的部分签名结果
func (service *htlcForwardTxManager) OnAliceSignedC3aAtAliceSide(msg bean.RequestMessage, user bean.User) (toAlice, toBob interface{}, err error) {

	if tool.CheckIsString(&msg.Data) == false {
		err = errors.New(enum.Tips_common_empty + "msg.data")
		log.Println(err)
		return nil, nil, err
	}

	signedDataForC3a := bean.AliceSignedHtlcDataForC3a{}
	_ = json.Unmarshal([]byte(msg.Data), &signedDataForC3a)

	if tool.CheckIsString(&signedDataForC3a.ChannelId) == false {
		err = errors.New(enum.Tips_common_empty + "channel_id")
		log.Println(err)
		return nil, nil, err
	}

	dataTo40P := service.tempDataSendTo40PAtAliceSide[user.PeerId+"_"+signedDataForC3a.ChannelId]
	if len(dataTo40P.ChannelId) == 0 {
		return nil, nil, errors.New(enum.Tips_common_wrong + "channel_id")
	}

	if tool.CheckIsString(&signedDataForC3a.C3aRsmcPartialSignedHex) {
		if pass, _ := rpcClient.CheckMultiSign(true, signedDataForC3a.C3aRsmcPartialSignedHex, 1); pass == false {
			err = errors.New(enum.Tips_common_wrong + "c3a_rsmc_partial_signed_hex")
			log.Println(err)
			return nil, nil, err
		}
	}

	if tool.CheckIsString(&signedDataForC3a.C3aHtlcPartialSignedHex) == false {
		err = errors.New(enum.Tips_common_empty + "c3a_htlc_partial_signed_hex")
		log.Println(err)
		return nil, nil, err
	}
	if pass, _ := rpcClient.CheckMultiSign(true, signedDataForC3a.C3aHtlcPartialSignedHex, 1); pass == false {
		err = errors.New(enum.Tips_common_wrong + "counterparty_signed_hex")
		log.Println(err)
		return nil, nil, err
	}

	if tool.CheckIsString(&signedDataForC3a.C3aCounterpartyPartialSignedHex) {
		if pass, _ := rpcClient.CheckMultiSign(true, signedDataForC3a.C3aCounterpartyPartialSignedHex, 1); pass == false {
			err = errors.New(enum.Tips_common_wrong + "c3a_counterparty_partial_signed_hex")
			log.Println(err)
			return nil, nil, err
		}
	}

	tx, err := user.Db.Begin(true)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	defer tx.Rollback()

	latestCommitmentTxInfo, err := getLatestCommitmentTxUseDbTx(tx, signedDataForC3a.ChannelId, user.PeerId)
	if err != nil {
		return nil, nil, errors.New(enum.Tips_channel_notFoundLatestCommitmentTx)
	}

	if len(latestCommitmentTxInfo.RSMCTxHex) > 0 {
		pass, _ := rpcClient.CheckMultiSign(true, signedDataForC3a.C3aRsmcPartialSignedHex, 1)
		if pass == false {
			return nil, nil, errors.New("error sign c3a_rsmc_partial_signed_hex")
		}
		//封装好的签名数据，给bob的客户端签名使用
		latestCommitmentTxInfo.RsmcRawTxData.Hex = signedDataForC3a.C3aRsmcPartialSignedHex
		latestCommitmentTxInfo.RSMCTxHex = signedDataForC3a.C3aRsmcPartialSignedHex
		latestCommitmentTxInfo.RSMCTxid = rpcClient.GetTxId(signedDataForC3a.C3aRsmcPartialSignedHex)
		latestCommitmentTxInfo.CurrState = dao.TxInfoState_Create
	}

	pass, _ := rpcClient.CheckMultiSign(true, signedDataForC3a.C3aHtlcPartialSignedHex, 1)
	if pass == false {
		return nil, nil, errors.New("error sign c3a_htlc_partial_signed_hex")
	}
	//封装好的签名数据，给bob的客户端签名使用
	latestCommitmentTxInfo.HtlcRawTxData.Hex = signedDataForC3a.C3aHtlcPartialSignedHex
	latestCommitmentTxInfo.HtlcTxHex = signedDataForC3a.C3aHtlcPartialSignedHex
	latestCommitmentTxInfo.HTLCTxid = rpcClient.GetTxId(signedDataForC3a.C3aHtlcPartialSignedHex)

	if len(latestCommitmentTxInfo.ToCounterpartyTxHex) > 0 {
		pass, _ := rpcClient.CheckMultiSign(true, signedDataForC3a.C3aCounterpartyPartialSignedHex, 1)
		if pass == false {
			return nil, nil, errors.New("error sign c3a_counterparty_partial_signed_hex")
		}
		//封装好的签名数据，给bob的客户端签名使用
		latestCommitmentTxInfo.ToCounterpartyRawTxData.Hex = signedDataForC3a.C3aCounterpartyPartialSignedHex
		latestCommitmentTxInfo.ToCounterpartyTxHex = signedDataForC3a.C3aCounterpartyPartialSignedHex
		latestCommitmentTxInfo.ToCounterpartyTxid = rpcClient.GetTxId(signedDataForC3a.C3aCounterpartyPartialSignedHex)
	}

	latestCommitmentTxInfo.CurrState = dao.TxInfoState_Create
	_ = tx.Update(latestCommitmentTxInfo)
	_ = tx.Commit()

	dataTo40P.C3aRsmcPartialSignedData = latestCommitmentTxInfo.RsmcRawTxData
	dataTo40P.C3aHtlcPartialSignedData = latestCommitmentTxInfo.HtlcRawTxData
	dataTo40P.C3aCounterpartyPartialSignedData = latestCommitmentTxInfo.ToCounterpartyRawTxData

	toAliceResult := bean.AliceSignedHtlcDataForC3aResult{}
	toAliceResult.ChannelId = dataTo40P.ChannelId
	toAliceResult.CommitmentTxHash = dataTo40P.PayerCommitmentTxHash
	return toAliceResult, dataTo40P, nil
}

// step 3 bob -40号协议 缓存来自40号协议的信息 推送110040消息，需要bob对C3a的交易进行签名
func (service *htlcForwardTxManager) BeforeBobSignAddHtlcRequestAtBobSide_40(msgData string, user bean.User) (data interface{}, err error) {
	requestAddHtlc := &bean.CreateHtlcTxForC3aOfP2p{}
	_ = json.Unmarshal([]byte(msgData), requestAddHtlc)
	channelId := requestAddHtlc.ChannelId

	tx, err := user.Db.Begin(true)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer tx.Rollback()

	channelInfo := &dao.ChannelInfo{}
	err = tx.Select(
		q.Eq("ChannelId", channelId),
		q.Or(
			q.Eq("PeerIdA", user.PeerId),
			q.Eq("PeerIdB", user.PeerId))).
		First(channelInfo)
	if channelInfo == nil {
		return nil, errors.New("not found channel info")
	}

	channelInfo.CurrState = dao.ChannelState_NewTx
	_ = tx.Update(channelInfo)

	_ = tx.Commit()

	if service.addHtlcTempDataAtBobSide == nil {
		service.addHtlcTempDataAtBobSide = make(map[string]string)
	}
	service.addHtlcTempDataAtBobSide[requestAddHtlc.PayerCommitmentTxHash] = msgData
	toBobData := bean.CreateHtlcTxForC3aToBob{}
	toBobData.ChannelId = requestAddHtlc.ChannelId
	toBobData.PayerCommitmentTxHash = requestAddHtlc.PayerCommitmentTxHash
	toBobData.PayerPeerId = requestAddHtlc.PayerPeerId
	toBobData.PayerNodeAddress = requestAddHtlc.PayerNodeAddress
	toBobData.C3aRsmcPartialSignedData = requestAddHtlc.C3aRsmcPartialSignedData
	toBobData.C3aCounterpartyPartialSignedData = requestAddHtlc.C3aCounterpartyPartialSignedData
	toBobData.C3aHtlcPartialSignedData = requestAddHtlc.C3aHtlcPartialSignedData
	return toBobData, nil
}

// step 4 bob 响应-100041号协议，创建C3a的Rsmc的Rd和Br，toHtlc的Br，Ht1a，Hlock，以及C3b的toB，toRsmc，toHtlc
func (service *htlcForwardTxManager) BobSignedAddHtlcAtBobSide(jsonData string, user bean.User) (returnData interface{}, err error) {
	if tool.CheckIsString(&jsonData) == false {
		err := errors.New(enum.Tips_common_empty + "msg data")
		log.Println(err)
		return nil, err
	}

	requestData := bean.BobSignedC3a{}
	err = json.Unmarshal([]byte(jsonData), &requestData)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	tx, err := user.Db.Begin(true)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer tx.Rollback()

	if tool.CheckIsString(&requestData.PayerCommitmentTxHash) == false {
		return nil, errors.New(enum.Tips_common_empty + "payer_commitment_tx_hash")
	}

	aliceMsg := service.addHtlcTempDataAtBobSide[requestData.PayerCommitmentTxHash]
	if tool.CheckIsString(&aliceMsg) == false {
		return nil, errors.New(enum.Tips_common_empty + "payer_commitment_tx_hash")
	}

	payerRequestAddHtlc := &bean.CreateHtlcTxForC3aOfP2p{}
	_ = json.Unmarshal([]byte(aliceMsg), payerRequestAddHtlc)

	if len(payerRequestAddHtlc.C3aRsmcPartialSignedData.Hex) > 0 {
		if pass, _ := rpcClient.CheckMultiSign(true, requestData.C3aCompleteSignedRsmcHex, 2); pass == false {
			err = errors.New(enum.Tips_common_empty + "c3a_complete_signed_rsmc_hex")
			log.Println(err)
			return nil, err
		}
		payerRequestAddHtlc.C3aRsmcPartialSignedData.Hex = requestData.C3aCompleteSignedRsmcHex
	}

	if pass, _ := rpcClient.CheckMultiSign(true, requestData.C3aCompleteSignedHtlcHex, 2); pass == false {
		err = errors.New(enum.Tips_common_empty + "c3a_complete_signed_htlc_hex")
		log.Println(err)
		return nil, err
	}
	payerRequestAddHtlc.C3aHtlcPartialSignedData.Hex = requestData.C3aCompleteSignedHtlcHex

	if len(payerRequestAddHtlc.C3aRsmcPartialSignedData.Hex) > 0 {
		if pass, _ := rpcClient.CheckMultiSign(true, requestData.C3aCompleteSignedCounterpartyHex, 2); pass == false {
			err = errors.New(enum.Tips_common_empty + "c3a_complete_signed_counterparty_hex")
			log.Println(err)
			return nil, err
		}
		payerRequestAddHtlc.C3aCounterpartyPartialSignedData.Hex = requestData.C3aCompleteSignedCounterpartyHex
	}

	channelId := payerRequestAddHtlc.ChannelId

	needBobSignData := bean.NeedBobSignHtlcTxOfC3b{}
	needBobSignData.ChannelId = channelId

	toAliceDataOf41P := bean.NeedAliceSignHtlcTxOfC3bP2p{}

	toAliceDataOf41P.PayerCommitmentTxHash = requestData.PayerCommitmentTxHash
	toAliceDataOf41P.PayeePeerId = user.PeerId
	toAliceDataOf41P.PayeeNodeAddress = user.P2PLocalPeerId

	if len(payerRequestAddHtlc.C3aRsmcPartialSignedData.Hex) > 0 {
		if pass, _ := rpcClient.CheckMultiSign(true, requestData.C3aCompleteSignedRsmcHex, 2); pass == false {
			return nil, errors.New(enum.Tips_common_wrong + "c3a_complete_signed_rsmc_hex")
		}
	}
	toAliceDataOf41P.C3aCompleteSignedRsmcHex = requestData.C3aCompleteSignedRsmcHex

	if pass, _ := rpcClient.CheckMultiSign(true, requestData.C3aCompleteSignedHtlcHex, 2); pass == false {
		return nil, errors.New(enum.Tips_common_wrong + "c3a_complete_signed_htlc_hex")
	}
	toAliceDataOf41P.C3aCompleteSignedHtlcHex = requestData.C3aCompleteSignedHtlcHex

	if len(payerRequestAddHtlc.C3aCounterpartyPartialSignedData.Hex) > 0 {
		if pass, _ := rpcClient.CheckMultiSign(true, requestData.C3aCompleteSignedCounterpartyHex, 2); pass == false {
			return nil, errors.New(enum.Tips_common_wrong + "c3a_complete_signed_counterparty_hex")
		}
	}
	toAliceDataOf41P.C3aCompleteSignedCounterpartyHex = requestData.C3aCompleteSignedCounterpartyHex

	// region check input data

	channelInfo := &dao.ChannelInfo{}
	err = tx.Select(
		q.Eq("ChannelId", channelId),
		q.Or(
			q.Eq("PeerIdA", user.PeerId),
			q.Eq("PeerIdB", user.PeerId))).
		First(channelInfo)
	if err != nil {
		return nil, errors.New(enum.Tips_htlc_noChanneFromRountingPacket)
	}
	err = checkBtcFundFinish(channelInfo.ChannelAddress, false)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	toAliceDataOf41P.ChannelId = channelInfo.ChannelId
	bobChannelPubKey := channelInfo.PubKeyB
	if user.PeerId == channelInfo.PeerIdA {
		bobChannelPubKey = channelInfo.PubKeyA
	}

	latestCommitmentTxInfo, _ := getLatestCommitmentTxUseDbTx(tx, channelInfo.ChannelId, user.PeerId)
	if latestCommitmentTxInfo.CurrState == dao.TxInfoState_Init {
		_ = tx.DeleteStruct(latestCommitmentTxInfo)
	} else {
		latestCommitmentTxInfo, _ = getLatestCommitmentTxUseDbTx(tx, channelInfo.ChannelId, user.PeerId)
	}

	if latestCommitmentTxInfo.Id > 0 {
		if tool.CheckIsString(&requestData.LastTempAddressPrivateKey) == false {
			err = errors.New(enum.Tips_common_empty + "last_temp_address_private_key")
			log.Println(err)
			return nil, err
		}

		if latestCommitmentTxInfo.CurrState == dao.TxInfoState_CreateAndSign {
			_, err = tool.GetPubKeyFromWifAndCheck(requestData.LastTempAddressPrivateKey, latestCommitmentTxInfo.RSMCTempAddressPubKey)
			if err != nil {
				return nil, errors.New(fmt.Sprintf(enum.Tips_rsmc_wrongPrivateKeyForLast, requestData.LastTempAddressPrivateKey, latestCommitmentTxInfo.RSMCTempAddressPubKey))
			}
		}
		if latestCommitmentTxInfo.CurrState == dao.TxInfoState_Create {
			if latestCommitmentTxInfo.TxType != dao.CommitmentTransactionType_Htlc {
				return nil, errors.New("error commitment tx type")
			}

			if requestData.CurrRsmcTempAddressPubKey != latestCommitmentTxInfo.RSMCTempAddressPubKey {
				return nil, errors.New(fmt.Sprintf(enum.Tips_rsmc_notSameValueWhenCreate, requestData.CurrRsmcTempAddressPubKey, latestCommitmentTxInfo.RSMCTempAddressPubKey))
			}
			if requestData.CurrHtlcTempAddressPubKey != latestCommitmentTxInfo.HTLCTempAddressPubKey {
				return nil, errors.New(fmt.Sprintf(enum.Tips_rsmc_notSameValueWhenCreate, requestData.CurrHtlcTempAddressPubKey, latestCommitmentTxInfo.HTLCTempAddressPubKey))
			}

			if latestCommitmentTxInfo.LastCommitmentTxId > 0 {
				lastCommitmentTx := &dao.CommitmentTransaction{}
				_ = tx.One("Id", latestCommitmentTxInfo.LastCommitmentTxId, lastCommitmentTx)
				_, err = tool.GetPubKeyFromWifAndCheck(requestData.LastTempAddressPrivateKey, lastCommitmentTx.RSMCTempAddressPubKey)
				if err != nil {
					return nil, err
				}
			}
		}
		toAliceDataOf41P.PayeeLastTempAddressPrivateKey = requestData.LastTempAddressPrivateKey
	}
	if tool.CheckIsString(&requestData.CurrRsmcTempAddressPubKey) == false {
		err = errors.New(enum.Tips_common_empty + "curr_rsmc_temp_address_pub_key")
		log.Println(err)
		return nil, err
	}
	toAliceDataOf41P.PayeeCurrRsmcTempAddressPubKey = requestData.CurrRsmcTempAddressPubKey

	if tool.CheckIsString(&requestData.CurrHtlcTempAddressPubKey) == false {
		err = errors.New(enum.Tips_common_empty + "curr_htlc_temp_address_pub_key")
		log.Println(err)
		return nil, err
	}
	toAliceDataOf41P.PayeeCurrHtlcTempAddressPubKey = requestData.CurrHtlcTempAddressPubKey
	//endregion

	//region 1、验证C3a的Rsmc的签名
	var c3aRsmcTxId, c3aSignedRsmcHex string
	var c3aRsmcMultiAddress, c3aRsmcRedeemScript, c3aRsmcMultiAddressScriptPubKey string
	var c3aRsmcOutputs []rpc.TransactionInputItem
	if tool.CheckIsString(&payerRequestAddHtlc.C3aRsmcPartialSignedData.Hex) {
		c3aSignedRsmcHex = requestData.C3aCompleteSignedRsmcHex
		testResult, err := rpcClient.TestMemPoolAccept(c3aSignedRsmcHex)
		if err != nil {
			return nil, err
		}
		if gjson.Parse(testResult).Array()[0].Get("allowed").Bool() == false {
			return nil, errors.New(gjson.Parse(testResult).Array()[0].Get("reject-reason").String())
		}
		c3aRsmcTxId = gjson.Parse(testResult).Array()[0].Get("txid").Str

		// region 根据alice的临时地址+bob的通道address,获取alice2+bob的多签地址，并得到AliceSignedRsmcHex签名后的交易的input，为创建alice的RD和bob的BR做准备
		c3aRsmcMultiAddress, c3aRsmcRedeemScript, c3aRsmcMultiAddressScriptPubKey, err = createMultiSig(payerRequestAddHtlc.CurrRsmcTempAddressPubKey, bobChannelPubKey)
		if err != nil {
			return nil, err
		}

		c3aRsmcOutputs, err = getInputsForNextTxByParseTxHashVout(c3aSignedRsmcHex, c3aRsmcMultiAddress, c3aRsmcMultiAddressScriptPubKey, c3aRsmcRedeemScript)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if len(c3aRsmcOutputs) == 0 {
			return nil, errors.New(enum.Tips_common_wrong + "payer rsmc hex")
		}
		//endregion
	}
	//endregion

	// region 2、验证C3a的toCounterpartyTxHex
	c3aSignedToCounterpartyTxHex := ""
	if len(payerRequestAddHtlc.C3aCounterpartyPartialSignedData.Hex) > 0 {
		c3aSignedToCounterpartyTxHex = payerRequestAddHtlc.C3aCounterpartyPartialSignedData.Hex
		testResult, err := rpcClient.TestMemPoolAccept(c3aSignedToCounterpartyTxHex)
		if err != nil {
			return nil, err
		}
		if gjson.Parse(testResult).Array()[0].Get("allowed").Bool() == false {
			return nil, errors.New(gjson.Parse(testResult).Array()[0].Get("reject-reason").String())
		}
	}
	//endregion

	// region 3、验证C3a的 htlcHex
	c3aSignedHtlcHex := requestData.C3aCompleteSignedHtlcHex
	testResult, err := rpcClient.TestMemPoolAccept(c3aSignedHtlcHex)
	if err != nil {
		return nil, err
	}
	c3aHtlcTxId := gjson.Parse(testResult).Array()[0].Get("txid").Str

	// region 根据alice的htlc临时地址+bob的通道address,获取alice2+bob的多签地址，并得到AliceSignedHtlcHex签名后的交易的input，为创建bob的HBR做准备
	c3aHtlcMultiAddress, c3aHtlcRedeemScript, c3aHtlcAddrScriptPubKey, err := createMultiSig(payerRequestAddHtlc.CurrHtlcTempAddressPubKey, bobChannelPubKey)
	if err != nil {
		return nil, err
	}

	c3aHtlcOutputs, err := getInputsForNextTxByParseTxHashVout(c3aSignedHtlcHex, c3aHtlcMultiAddress, c3aHtlcAddrScriptPubKey, c3aHtlcRedeemScript)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//endregion
	//endregion

	//获取bob最新的承诺交易
	isFirstRequest := false
	if latestCommitmentTxInfo != nil && latestCommitmentTxInfo.Id > 0 {
		if latestCommitmentTxInfo.CurrState == dao.TxInfoState_CreateAndSign {
			if latestCommitmentTxInfo.TxType != dao.CommitmentTransactionType_Rsmc {
				return nil, errors.New("wrong commitment tx type " + strconv.Itoa(int(latestCommitmentTxInfo.TxType)))
			}
		}
		if latestCommitmentTxInfo.CurrState != dao.TxInfoState_CreateAndSign && latestCommitmentTxInfo.CurrState != dao.TxInfoState_Create {
			return nil, errors.New("wrong commitment tx state " + strconv.Itoa(int(latestCommitmentTxInfo.CurrState)))
		}

		if latestCommitmentTxInfo.CurrState == dao.TxInfoState_CreateAndSign { //有上一次的承诺交易
			isFirstRequest = true
		}
	} else { // 因为没有充值，没有最初的承诺交易C1b
		isFirstRequest = true
	}

	var amountToOther = 0.0
	var amountToHtlc = 0.0
	//如果是本轮的第一次请求交易
	if isFirstRequest {
		//region 4、根据对方传过来的上一个交易的临时rsmc私钥，签名最近的BR交易，保证对方确实放弃了上一个承诺交易
		err := signLastBR(tx, dao.BRType_Rmsc, *channelInfo, user.PeerId, payerRequestAddHtlc.LastTempAddressPrivateKey, latestCommitmentTxInfo.Id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		//endregion

		fundingTransaction := getFundingTransactionByChannelId(tx, channelInfo.ChannelId, user.PeerId)
		if fundingTransaction == nil {
			return nil, errors.New(enum.Tips_common_notFound + "fundingTransaction")
		}

		//region 5、创建C3b
		latestCommitmentTxInfo, err = htlcPayeeCreateCommitmentTx_C3b(tx, channelInfo, requestData, *payerRequestAddHtlc, latestCommitmentTxInfo, c3aSignedToCounterpartyTxHex, user)
		//endregion
	}

	amountToOther = latestCommitmentTxInfo.AmountToCounterparty
	amountToHtlc = latestCommitmentTxInfo.AmountToHtlc

	needBobSignData.C3bCounterpartyRawData = latestCommitmentTxInfo.ToCounterpartyRawTxData
	needBobSignData.C3bRsmcRawData = latestCommitmentTxInfo.RsmcRawTxData
	needBobSignData.C3bHtlcRawData = latestCommitmentTxInfo.HtlcRawTxData

	toAliceDataOf41P.C3bCounterpartyPartialSignedData = latestCommitmentTxInfo.ToCounterpartyRawTxData
	toAliceDataOf41P.C3bRsmcPartialSignedData = latestCommitmentTxInfo.RsmcRawTxData
	toAliceDataOf41P.C3bHtlcPartialSignedData = latestCommitmentTxInfo.HtlcRawTxData
	toAliceDataOf41P.PayeeCommitmentTxHash = latestCommitmentTxInfo.CurrHash

	var myAddress = channelInfo.AddressB
	if user.PeerId == channelInfo.PeerIdA {
		myAddress = channelInfo.AddressA
	}
	tempC3aSideCommitmentTx := &dao.CommitmentTransaction{}

	if len(c3aRsmcOutputs) > 0 {

		//region 6、根据alice C3a的Rsmc输出，创建对应的BR,为下一个交易做准备，create BR2b tx  for bob
		tempC3aSideCommitmentTx.Id = latestCommitmentTxInfo.Id
		tempC3aSideCommitmentTx.PropertyId = channelInfo.PropertyId
		tempC3aSideCommitmentTx.RSMCTempAddressPubKey = payerRequestAddHtlc.CurrRsmcTempAddressPubKey
		tempC3aSideCommitmentTx.RSMCMultiAddress = c3aRsmcMultiAddress
		tempC3aSideCommitmentTx.RSMCMultiAddressScriptPubKey = c3aRsmcMultiAddressScriptPubKey
		tempC3aSideCommitmentTx.RSMCRedeemScript = c3aRsmcRedeemScript
		tempC3aSideCommitmentTx.RSMCTxHex = c3aSignedRsmcHex
		tempC3aSideCommitmentTx.RSMCTxid = c3aRsmcTxId
		tempC3aSideCommitmentTx.AmountToRSMC = latestCommitmentTxInfo.AmountToCounterparty

		c3aRsmcBrHexData, err := createRawBR(dao.BRType_Rmsc, channelInfo, tempC3aSideCommitmentTx, c3aRsmcOutputs, myAddress, user)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		c3aRsmcBrHexData.PubKeyA = payerRequestAddHtlc.CurrRsmcTempAddressPubKey
		c3aRsmcBrHexData.PubKeyB = bobChannelPubKey
		needBobSignData.C3aRsmcBrRawData = c3aRsmcBrHexData
		//endregion

		//region 7、根据签名后的AliceRsmc创建alice的RD create RD tx for alice
		aliceRdOutputAddress := channelInfo.AddressA
		if user.PeerId == channelInfo.PeerIdA {
			aliceRdOutputAddress = channelInfo.AddressB
		}
		c3aRsmcRdData, err := rpcClient.OmniCreateRawTransactionUseUnsendInput(
			c3aRsmcMultiAddress,
			c3aRsmcOutputs,
			aliceRdOutputAddress,
			channelInfo.FundingAddress,
			channelInfo.PropertyId,
			amountToOther,
			getBtcMinerAmount(channelInfo.BtcAmount),
			1000,
			&c3aRsmcRedeemScript)
		if err != nil {
			log.Println(err)
			return nil, errors.New(fmt.Sprintf(enum.Tips_rsmc_failToCreate, "RD raw transacation"))
		}

		signHexData := bean.NeedClientSignTxData{}
		signHexData.Hex = c3aRsmcRdData["hex"].(string)
		signHexData.Inputs = c3aRsmcRdData["inputs"]
		signHexData.IsMultisig = true
		signHexData.PubKeyA = payerRequestAddHtlc.CurrRsmcTempAddressPubKey
		signHexData.PubKeyB = bobChannelPubKey
		needBobSignData.C3aRsmcRdRawData = signHexData
		toAliceDataOf41P.C3aRsmcRdPartialSignedData = signHexData
		//endregion
	}

	// region 7、根据alice C3a的Htlc输出，创建对应的BR,为下一个交易做准备，create HBR2b tx  for bob
	tempC3aSideCommitmentTx.Id = latestCommitmentTxInfo.Id
	tempC3aSideCommitmentTx.PropertyId = channelInfo.PropertyId
	tempC3aSideCommitmentTx.RSMCTempAddressPubKey = payerRequestAddHtlc.CurrHtlcTempAddressPubKey
	tempC3aSideCommitmentTx.RSMCMultiAddress = c3aHtlcMultiAddress
	tempC3aSideCommitmentTx.RSMCRedeemScript = c3aHtlcRedeemScript
	tempC3aSideCommitmentTx.RSMCMultiAddressScriptPubKey = c3aHtlcAddrScriptPubKey
	tempC3aSideCommitmentTx.RSMCTxHex = c3aSignedHtlcHex
	tempC3aSideCommitmentTx.RSMCTxid = c3aHtlcTxId
	tempC3aSideCommitmentTx.AmountToRSMC = latestCommitmentTxInfo.AmountToHtlc
	c3aHtlcBrHexData, err := createRawBR(dao.BRType_Htlc, channelInfo, tempC3aSideCommitmentTx, c3aHtlcOutputs, myAddress, user)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	c3aHtlcBrHexData.PubKeyA = payerRequestAddHtlc.CurrHtlcTempAddressPubKey
	c3aHtlcBrHexData.PubKeyB = bobChannelPubKey
	needBobSignData.C3aHtlcBrRawData = c3aHtlcBrHexData
	//endregion

	// region  8、h+bobChannelPubkey 锁定给bob的付款金额
	lockByHForBobTx, err := createHtlcLockByHForBobAtPayeeSide(*channelInfo, *payerRequestAddHtlc, c3aSignedHtlcHex, bobChannelPubKey, channelInfo.PropertyId, amountToHtlc)
	if err != nil {
		return nil, err
	}
	lockByHForBobTx.PubKeyA = payerRequestAddHtlc.CurrHtlcTempAddressPubKey
	lockByHForBobTx.PubKeyB = bobChannelPubKey

	needBobSignData.C3aHtlcHlockRawData = *lockByHForBobTx
	toAliceDataOf41P.C3aHtlcHlockPartialSignedData = *lockByHForBobTx
	//endregion

	// region 9、ht1a 根据signedHtlcHex（alice签名后C3a的第三个输出）作为输入生成
	ht1aTxData, err := createHT1aForAlice(*channelInfo, *payerRequestAddHtlc, c3aSignedHtlcHex, bobChannelPubKey, channelInfo.PropertyId, amountToHtlc, latestCommitmentTxInfo.HtlcCltvExpiry)
	if err != nil {
		return nil, err
	}
	ht1aTxData.PubKeyA = bobChannelPubKey
	ht1aTxData.PubKeyB = payerRequestAddHtlc.CurrHtlcTempAddressPubKey
	needBobSignData.C3aHtlcHtRawData = *ht1aTxData
	toAliceDataOf41P.C3aHtlcHtPartialSignedData = *ht1aTxData
	toAliceDataOf41P.C3aHtlcTempAddressForHtPubKey = payerRequestAddHtlc.CurrHtlcTempAddressForHt1aPubKey
	//endregion

	channelInfo.CurrState = dao.ChannelState_HtlcTx
	_ = tx.Update(channelInfo)

	_ = tx.Commit()

	if service.tempDataSendTo41PAtBobSide == nil {
		service.tempDataSendTo41PAtBobSide = make(map[string]bean.NeedAliceSignHtlcTxOfC3bP2p)
	}
	service.tempDataSendTo41PAtBobSide[user.PeerId+"_"+channelInfo.ChannelId] = toAliceDataOf41P

	return needBobSignData, nil
}

// step 5 bob -100101 bob完成对C3b的签名，构建41号协议的消息体，推送41号协议
func (service *htlcForwardTxManager) OnBobSignedC3bAtBobSide(msg bean.RequestMessage, user bean.User) (toAlice, toBob interface{}, err error) {
	c3bResult := bean.BobSignedHtlcTxOfC3b{}
	err = json.Unmarshal([]byte(msg.Data), &c3bResult)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	key := user.PeerId + "_" + c3bResult.ChannelId
	toAliceDataOfP2p := service.tempDataSendTo41PAtBobSide[key]

	if len(toAliceDataOfP2p.ChannelId) == 0 {
		return nil, nil, errors.New(enum.Tips_common_wrong + "channel_id")
	}

	aliceMsg := service.addHtlcTempDataAtBobSide[toAliceDataOfP2p.PayerCommitmentTxHash]
	if tool.CheckIsString(&aliceMsg) == false {
		return nil, nil, errors.New(enum.Tips_common_empty + "payer_commitment_tx_hash")
	}
	payerRequestAddHtlc := &bean.CreateHtlcTxForC3aOfP2p{}
	_ = json.Unmarshal([]byte(aliceMsg), payerRequestAddHtlc)

	if len(payerRequestAddHtlc.C3aRsmcPartialSignedData.Hex) > 0 {
		if pass, _ := rpcClient.CheckMultiSign(true, toAliceDataOfP2p.C3aCompleteSignedRsmcHex, 2); pass == false {
			return nil, nil, errors.New(enum.Tips_common_wrong + "c3a_complete_signed_rsmc_hex")
		}
	}

	if pass, _ := rpcClient.CheckMultiSign(true, toAliceDataOfP2p.C3aCompleteSignedHtlcHex, 2); pass == false {
		return nil, nil, errors.New(enum.Tips_common_wrong + "c3a_complete_signed_htlc_hex")
	}

	if len(payerRequestAddHtlc.C3aCounterpartyPartialSignedData.Hex) > 0 {
		if pass, _ := rpcClient.CheckMultiSign(true, toAliceDataOfP2p.C3aCompleteSignedCounterpartyHex, 2); pass == false {
			return nil, nil, errors.New(enum.Tips_common_wrong + "c3a_complete_signed_counterparty_hex")
		}
	}

	tx, err := user.Db.Begin(true)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	defer tx.Rollback()

	channelInfo := &dao.ChannelInfo{}
	err = tx.Select(
		q.Eq("ChannelId", c3bResult.ChannelId),
		q.Or(
			q.Eq("PeerIdA", user.PeerId),
			q.Eq("PeerIdB", user.PeerId))).
		First(channelInfo)
	if err != nil {
		return nil, nil, errors.New(enum.Tips_htlc_noChanneFromRountingPacket)
	}

	latestCommitmentTxInfo, _ := getLatestCommitmentTxUseDbTx(tx, channelInfo.ChannelId, user.PeerId)
	if latestCommitmentTxInfo.CurrState != dao.TxInfoState_Init {
		return nil, nil, errors.New(enum.Tips_channel_notFoundLatestCommitmentTx)
	}

	tempC3aSideCommitmentTx := &dao.CommitmentTransaction{}
	bobChannelPubKey := channelInfo.PubKeyB
	var myAddress = channelInfo.AddressB
	if user.PeerId == channelInfo.PeerIdA {
		myAddress = channelInfo.AddressA
		bobChannelPubKey = channelInfo.PubKeyA
	}

	if tool.CheckIsString(&toAliceDataOfP2p.C3aRsmcRdPartialSignedData.Hex) {
		if pass, _ := rpcClient.CheckMultiSign(false, c3bResult.C3aRsmcRdPartialSignedHex, 1); pass == false {
			return nil, nil, errors.New(enum.Tips_common_wrong + "c3a_rsmc_rd_partial_signed_hex")
		}
		toAliceDataOfP2p.C3aRsmcRdPartialSignedData.Hex = c3bResult.C3aRsmcRdPartialSignedHex

		if pass, _ := rpcClient.CheckMultiSign(false, c3bResult.C3aRsmcBrPartialSignedHex, 1); pass == false {
			return nil, nil, errors.New(enum.Tips_common_wrong + "c3a_rsmc_br_partial_signed_hex")
		}

		c3aSignedRsmcHex := toAliceDataOfP2p.C3aCompleteSignedRsmcHex
		c3aRsmcTxId := rpcClient.GetTxId(c3aSignedRsmcHex)

		c3aRsmcMultiAddress, c3aRsmcRedeemScript, c3aRsmcMultiAddressScriptPubKey, err := createMultiSig(payerRequestAddHtlc.CurrRsmcTempAddressPubKey, bobChannelPubKey)
		if err != nil {
			return nil, nil, err
		}

		c3aRsmcOutputs, err := getInputsForNextTxByParseTxHashVout(c3aSignedRsmcHex, c3aRsmcMultiAddress, c3aRsmcMultiAddressScriptPubKey, c3aRsmcRedeemScript)
		if err != nil {
			log.Println(err)
			return nil, nil, err
		}

		tempC3aSideCommitmentTx.Id = latestCommitmentTxInfo.Id
		tempC3aSideCommitmentTx.PropertyId = channelInfo.PropertyId
		tempC3aSideCommitmentTx.RSMCTempAddressPubKey = payerRequestAddHtlc.CurrRsmcTempAddressPubKey
		tempC3aSideCommitmentTx.RSMCMultiAddress = c3aRsmcMultiAddress
		tempC3aSideCommitmentTx.RSMCMultiAddressScriptPubKey = c3aRsmcMultiAddressScriptPubKey
		tempC3aSideCommitmentTx.RSMCRedeemScript = c3aRsmcRedeemScript
		tempC3aSideCommitmentTx.RSMCTxHex = toAliceDataOfP2p.C3aCompleteSignedRsmcHex
		tempC3aSideCommitmentTx.RSMCTxid = c3aRsmcTxId
		tempC3aSideCommitmentTx.AmountToRSMC = latestCommitmentTxInfo.AmountToCounterparty
		_ = createCurrCommitmentTxPartialSignedBR(tx, dao.BRType_Rmsc, channelInfo, tempC3aSideCommitmentTx, c3aRsmcOutputs, myAddress, c3bResult.C3aRsmcBrPartialSignedHex, user)
	}

	if pass, _ := rpcClient.CheckMultiSign(false, c3bResult.C3aHtlcHtPartialSignedHex, 1); pass == false {
		return nil, nil, errors.New(enum.Tips_common_wrong + "c3a_htlc_ht_partial_signed_hex")
	}
	toAliceDataOfP2p.C3aHtlcHtPartialSignedData.Hex = c3bResult.C3aHtlcHtPartialSignedHex

	if pass, _ := rpcClient.CheckMultiSign(false, c3bResult.C3aHtlcHlockPartialSignedHex, 1); pass == false {
		return nil, nil, errors.New(enum.Tips_common_wrong + "c3a_htlc_hlock_partial_signed_hex")
	}
	toAliceDataOfP2p.C3aHtlcHlockPartialSignedData.Hex = c3bResult.C3aHtlcHlockPartialSignedHex

	if pass, _ := rpcClient.CheckMultiSign(false, c3bResult.C3aHtlcBrPartialSignedHex, 1); pass == false {
		return nil, nil, errors.New(enum.Tips_common_wrong + "c3a_htlc_br_partial_signed_hex")
	}

	c3aSignedHtlcHex := toAliceDataOfP2p.C3aCompleteSignedHtlcHex
	c3aHtlcMultiAddress, c3aHtlcRedeemScript, c3aHtlcAddrScriptPubKey, err := createMultiSig(payerRequestAddHtlc.CurrHtlcTempAddressPubKey, bobChannelPubKey)
	if err != nil {
		return nil, nil, err
	}

	c3aHtlcOutputs, err := getInputsForNextTxByParseTxHashVout(c3aSignedHtlcHex, c3aHtlcMultiAddress, c3aHtlcAddrScriptPubKey, c3aHtlcRedeemScript)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	tempC3aSideCommitmentTx.Id = latestCommitmentTxInfo.Id
	tempC3aSideCommitmentTx.PropertyId = channelInfo.PropertyId
	tempC3aSideCommitmentTx.RSMCTempAddressPubKey = payerRequestAddHtlc.CurrHtlcTempAddressPubKey
	tempC3aSideCommitmentTx.RSMCMultiAddress = c3aHtlcMultiAddress
	tempC3aSideCommitmentTx.RSMCRedeemScript = c3aHtlcRedeemScript
	tempC3aSideCommitmentTx.RSMCMultiAddressScriptPubKey = c3aHtlcAddrScriptPubKey
	tempC3aSideCommitmentTx.RSMCTxHex = c3aSignedHtlcHex
	tempC3aSideCommitmentTx.RSMCTxid = rpcClient.GetTxId(c3aSignedHtlcHex)
	tempC3aSideCommitmentTx.AmountToRSMC = latestCommitmentTxInfo.AmountToHtlc
	_ = createCurrCommitmentTxPartialSignedBR(tx, dao.BRType_Htlc, channelInfo, tempC3aSideCommitmentTx, c3aHtlcOutputs, myAddress, c3bResult.C3aHtlcBrPartialSignedHex, user)

	if tool.CheckIsString(&toAliceDataOfP2p.C3bRsmcPartialSignedData.Hex) {
		if pass, _ := rpcClient.CheckMultiSign(true, c3bResult.C3bRsmcPartialSignedHex, 1); pass == false {
			return nil, nil, errors.New(enum.Tips_common_wrong + "c3b_rsmc_partial_signed_hex")
		}
		toAliceDataOfP2p.C3bRsmcPartialSignedData.Hex = c3bResult.C3bRsmcPartialSignedHex
	}
	if tool.CheckIsString(&toAliceDataOfP2p.C3bCounterpartyPartialSignedData.Hex) {
		if pass, _ := rpcClient.CheckMultiSign(true, c3bResult.C3bCounterpartyPartialSignedHex, 1); pass == false {
			return nil, nil, errors.New(enum.Tips_common_wrong + "c3b_counterparty_partial_signed_hex")
		}
		toAliceDataOfP2p.C3bCounterpartyPartialSignedData.Hex = c3bResult.C3bCounterpartyPartialSignedHex
	}

	if pass, _ := rpcClient.CheckMultiSign(true, c3bResult.C3bHtlcPartialSignedHex, 1); pass == false {
		return nil, nil, errors.New(enum.Tips_common_wrong + "c3b_htlc_partial_signed_hex")
	}
	toAliceDataOfP2p.C3bHtlcPartialSignedData.Hex = c3bResult.C3bHtlcPartialSignedHex

	latestCommitmentTxInfo.CurrState = dao.TxInfoState_Create
	tx.UpdateField(latestCommitmentTxInfo, "CurrState", dao.TxInfoState_Create)
	tx.Commit()

	delete(service.tempDataSendTo41PAtBobSide, key)
	delete(service.addHtlcTempDataAtBobSide, payerRequestAddHtlc.PayerCommitmentTxHash)

	toBobData := bean.BobSignedHtlcTxOfC3bResult{}
	toBobData.ChannelId = channelInfo.ChannelId
	toBobData.CommitmentTxHash = latestCommitmentTxInfo.CurrHash
	return toAliceDataOfP2p, toBobData, nil
}

// step 6 alice p2p 41号协议，构建需要alice签名的数据，缓存41号协议的数据， 推送（110041）信息给alice签名
func (service *htlcForwardTxManager) AfterBobSignAddHtlcAtAliceSide_41(msgData string, user bean.User) (data interface{}, needNoticeBob bool, err error) {
	dataFromBob := bean.NeedAliceSignHtlcTxOfC3bP2p{}
	_ = json.Unmarshal([]byte(msgData), &dataFromBob)

	channelId := dataFromBob.ChannelId
	commitmentTxHash := dataFromBob.PayerCommitmentTxHash

	tx, err := user.Db.Begin(true)
	if err != nil {
		log.Println(err)
		return nil, false, err
	}
	defer tx.Rollback()
	commitmentTransaction := &dao.CommitmentTransaction{}
	err = tx.Select(
		q.Eq("CurrHash", commitmentTxHash),
		q.Eq("ChannelId", channelId)).First(commitmentTransaction)
	if err != nil {
		return nil, true, err
	}

	channelInfo := &dao.ChannelInfo{}
	err = tx.Select(
		q.Eq("ChannelId", channelId),
		q.Or(
			q.Eq("PeerIdA", user.PeerId),
			q.Eq("PeerIdB", user.PeerId))).
		First(channelInfo)
	if channelInfo == nil {
		return nil, true, errors.New("not found the channel " + channelId)
	}

	htlcRequestInfo := &dao.AddHtlcRequestInfo{}
	err = tx.Select(
		q.Eq("ChannelId", channelId),
		q.Eq("H", commitmentTransaction.HtlcH),
		q.Eq("PropertyId", channelInfo.PropertyId)).
		OrderBy("CreateAt").
		Reverse().
		First(htlcRequestInfo)
	if err != nil {
		return nil, false, err
	}

	needAliceSignHtlcTxOfC3b := bean.NeedAliceSignHtlcTxOfC3b{}
	needAliceSignHtlcTxOfC3b.ChannelId = channelId
	needAliceSignHtlcTxOfC3b.C3aRsmcRdPartialSignedData = dataFromBob.C3aRsmcRdPartialSignedData
	needAliceSignHtlcTxOfC3b.C3aHtlcHlockPartialSignedData = dataFromBob.C3aHtlcHlockPartialSignedData
	needAliceSignHtlcTxOfC3b.C3aHtlcHtPartialSignedData = dataFromBob.C3aHtlcHtPartialSignedData
	needAliceSignHtlcTxOfC3b.C3bRsmcPartialSignedData = dataFromBob.C3bRsmcPartialSignedData
	needAliceSignHtlcTxOfC3b.C3bHtlcPartialSignedData = dataFromBob.C3bHtlcPartialSignedData
	needAliceSignHtlcTxOfC3b.C3bCounterpartyPartialSignedData = dataFromBob.C3bCounterpartyPartialSignedData

	if service.tempDataFrom41PAtAliceSide == nil {
		service.tempDataFrom41PAtAliceSide = make(map[string]bean.NeedAliceSignHtlcTxOfC3bP2p)
	}
	service.tempDataFrom41PAtAliceSide[user.PeerId+"_"+channelId] = dataFromBob
	return needAliceSignHtlcTxOfC3b, false, nil
}

// step 7 alice 响应100102号协议，根据C3b的签名结果创建C3b的rmsc的rd，br，htlc的br，htd，hlock，以及创建C3a的htrd和htbr
func (service *htlcForwardTxManager) OnAliceSignC3bAtAliceSide(msg bean.RequestMessage, user bean.User) (interface{}, error) {

	aliceSignedC3b := bean.AliceSignedHtlcTxOfC3bResult{}
	_ = json.Unmarshal([]byte(msg.Data), &aliceSignedC3b)

	channelId := aliceSignedC3b.ChannelId
	dataFromBob := service.tempDataFrom41PAtAliceSide[user.PeerId+"_"+channelId]

	//为了准备给42传数据
	needBobSignData := bean.NeedBobSignHtlcSubTxOfC3bP2p{}
	needBobSignData.ChannelId = channelId
	needBobSignData.PayeeCommitmentTxHash = dataFromBob.PayeeCommitmentTxHash
	needBobSignData.C3aHtlcTempAddressForHtPubKey = dataFromBob.C3aHtlcTempAddressForHtPubKey
	needBobSignData.PayerPeerId = user.PeerId
	needBobSignData.PayerNodeAddress = user.P2PLocalPeerId

	needAliceSignData := bean.NeedAliceSignHtlcSubTxOfC3b{}
	needAliceSignData.ChannelId = channelId
	needAliceSignData.PayeePeerId = dataFromBob.PayeePeerId
	needAliceSignData.PayeeNodeAddress = dataFromBob.PayeeNodeAddress

	tx, err := user.Db.Begin(true)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	commitmentTransaction := &dao.CommitmentTransaction{}
	err = tx.Select(
		q.Eq("CurrHash", dataFromBob.PayerCommitmentTxHash),
		q.Eq("ChannelId", channelId)).First(commitmentTransaction)
	if err != nil {
		return nil, err
	}

	channelInfo := &dao.ChannelInfo{}
	err = tx.Select(
		q.Eq("ChannelId", channelId),
		q.Or(
			q.Eq("PeerIdA", user.PeerId),
			q.Eq("PeerIdB", user.PeerId))).
		First(channelInfo)
	if channelInfo == nil {
		return nil, errors.New("not found the channel " + channelId)
	}

	bobRsmcHexIsExist := false
	var c3bSignedRsmcHex, c3bSignedRsmcTxid, c3bSignedToCounterpartyTxHex string
	if len(dataFromBob.C3bRsmcPartialSignedData.Hex) > 0 {
		if pass, _ := rpcClient.CheckMultiSign(true, aliceSignedC3b.C3bRsmcCompleteSignedHex, 2); pass == false {
			return nil, errors.New(enum.Tips_common_wrong + "c3b_rsmc_complete_signed_hex")
		}
		bobRsmcHexIsExist = true
		c3bSignedRsmcHex = aliceSignedC3b.C3bRsmcCompleteSignedHex
		dataFromBob.C3bRsmcPartialSignedData.Hex = c3bSignedRsmcHex
		c3bSignedRsmcTxid = rpcClient.GetTxId(c3bSignedRsmcHex)
	}

	if len(dataFromBob.C3bCounterpartyPartialSignedData.Hex) > 0 {
		if pass, _ := rpcClient.CheckMultiSign(true, aliceSignedC3b.C3bCounterpartyCompleteSignedHex, 2); pass == false {
			return nil, errors.New(enum.Tips_common_wrong + "c3b_counterparty_complete_signed_hex")
		}
		c3bSignedToCounterpartyTxHex = aliceSignedC3b.C3bCounterpartyCompleteSignedHex
		dataFromBob.C3bCounterpartyPartialSignedData.Hex = c3bSignedToCounterpartyTxHex
	}

	if pass, _ := rpcClient.CheckMultiSign(true, aliceSignedC3b.C3bHtlcCompleteSignedHex, 2); pass == false {
		return nil, errors.New(enum.Tips_common_wrong + "c3b_htlc_complete_signed_hex")
	}
	c3bSignedHtlcHex := aliceSignedC3b.C3bHtlcCompleteSignedHex
	dataFromBob.C3bHtlcPartialSignedData.Hex = c3bSignedHtlcHex

	if pass, _ := rpcClient.CheckMultiSign(false, aliceSignedC3b.C3aHtlcHtCompleteSignedHex, 2); pass == false {
		return nil, errors.New(enum.Tips_common_wrong + "c3a_htlc_ht_complete_signed_hex")
	}

	if pass, _ := rpcClient.CheckMultiSign(false, aliceSignedC3b.C3aHtlcHlockCompleteSignedHex, 2); pass == false {
		return nil, errors.New(enum.Tips_common_wrong + "c3a_htlc_hlock_complete_signed_hex")
	}

	if pass, _ := rpcClient.CheckMultiSign(false, aliceSignedC3b.C3aRsmcRdCompleteSignedHex, 2); pass == false {
		return nil, errors.New(enum.Tips_common_wrong + "c3a_rsmc_rd_complete_signed_hex")
	}

	payerChannelPubKey := channelInfo.PubKeyA
	payerChannelAddress := channelInfo.AddressA
	payeeChannelAddress := channelInfo.AddressB
	payeeChannelPubKey := channelInfo.PubKeyB
	if user.PeerId == channelInfo.PeerIdB {
		payerChannelPubKey = channelInfo.PubKeyB
		payerChannelAddress = channelInfo.AddressB
		payeeChannelAddress = channelInfo.AddressA
		payeeChannelPubKey = channelInfo.PubKeyA
	}

	//region 1 c3b rsmcHex
	needBobSignData.C3bCompleteSignedRsmcHex = c3bSignedRsmcHex
	commitmentTransaction.FromCounterpartySideForMeTxHex = c3bSignedRsmcHex
	//endregion

	// region 2 c3b toCounterpartyTxHex
	needBobSignData.C3bCompleteSignedCounterpartyHex = c3bSignedToCounterpartyTxHex
	//endregion

	// region 3 c3b htlcHex
	needBobSignData.C3bCompleteSignedHtlcHex = c3bSignedHtlcHex
	//endregion

	var bobRsmcOutputs []rpc.TransactionInputItem
	if bobRsmcHexIsExist {

		//region 4 c3b Rsmc rd
		c3bRsmcMultiAddress, c3bRsmcRedeemScript, c3bRsmcAddrScriptPubKey, err := createMultiSig(dataFromBob.PayeeCurrRsmcTempAddressPubKey, payerChannelPubKey)
		if err != nil {
			return nil, err
		}
		bobRsmcOutputs, err = getInputsForNextTxByParseTxHashVout(c3bSignedRsmcHex, c3bRsmcMultiAddress, c3bRsmcAddrScriptPubKey, c3bRsmcRedeemScript)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		c3bRsmcRdTx, err := rpcClient.OmniCreateRawTransactionUseUnsendInput(
			c3bRsmcMultiAddress,
			bobRsmcOutputs,
			payeeChannelAddress,
			channelInfo.FundingAddress,
			channelInfo.PropertyId,
			commitmentTransaction.AmountToCounterparty,
			getBtcMinerAmount(channelInfo.BtcAmount),
			1000,
			&c3bRsmcRedeemScript)
		if err != nil {
			log.Println(err)
			return nil, errors.New("fail to create rd for c3b rsmc")
		}
		rdRawData := bean.NeedClientSignTxData{}
		rdRawData.Hex = c3bRsmcRdTx["hex"].(string)
		rdRawData.Inputs = c3bRsmcRdTx["inputs"]
		rdRawData.IsMultisig = true
		rdRawData.PubKeyA = dataFromBob.PayeeCurrRsmcTempAddressPubKey
		rdRawData.PubKeyB = payerChannelPubKey
		needBobSignData.C3bRsmcRdPartialData = rdRawData
		needAliceSignData.C3bRsmcRdRawData = rdRawData
		//endregion create rd tx for alice

		//region 5 c3b Rsmc br
		tempOtherSideCommitmentTx := &dao.CommitmentTransaction{}
		tempOtherSideCommitmentTx.Id = commitmentTransaction.Id
		tempOtherSideCommitmentTx.PropertyId = channelInfo.PropertyId
		tempOtherSideCommitmentTx.RSMCTempAddressPubKey = dataFromBob.PayeeCurrRsmcTempAddressPubKey
		tempOtherSideCommitmentTx.RSMCMultiAddress = c3bRsmcMultiAddress
		tempOtherSideCommitmentTx.RSMCRedeemScript = c3bRsmcRedeemScript
		tempOtherSideCommitmentTx.RSMCMultiAddressScriptPubKey = c3bRsmcAddrScriptPubKey
		tempOtherSideCommitmentTx.RSMCTxHex = c3bSignedRsmcHex
		tempOtherSideCommitmentTx.RSMCTxid = c3bSignedRsmcTxid
		tempOtherSideCommitmentTx.AmountToRSMC = commitmentTransaction.AmountToCounterparty
		rawBR, err := createRawBR(dao.BRType_Rmsc, channelInfo, tempOtherSideCommitmentTx, bobRsmcOutputs, payerChannelAddress, user)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		rawBR.PubKeyA = dataFromBob.PayeeCurrRsmcTempAddressPubKey
		rawBR.PubKeyB = payerChannelPubKey
		needAliceSignData.C3bRsmcBrRawData = rawBR
		//endregion create br tx for alice
	}

	//region 6 c3b htlc htd
	htlcTimeOut := commitmentTransaction.HtlcCltvExpiry
	bobHtlcMultiAddress, bobHtlcRedeemScript, bobHtlcMultiAddressScriptPubKey, err := createMultiSig(dataFromBob.PayeeCurrHtlcTempAddressPubKey, payerChannelPubKey)
	if err != nil {
		return nil, err
	}
	bobHtlcOutputs, err := getInputsForNextTxByParseTxHashVout(c3bSignedHtlcHex, bobHtlcMultiAddress, bobHtlcMultiAddressScriptPubKey, bobHtlcRedeemScript)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	c3bHtdTx, err := rpcClient.OmniCreateRawTransactionUseUnsendInput(
		bobHtlcMultiAddress,
		bobHtlcOutputs,
		payerChannelAddress,
		channelInfo.FundingAddress,
		channelInfo.PropertyId,
		commitmentTransaction.AmountToHtlc,
		getBtcMinerAmount(channelInfo.BtcAmount),
		htlcTimeOut,
		&bobHtlcRedeemScript)
	if err != nil {
		log.Println(err)
		return nil, errors.New("fail to create HTD1b for C3b")
	}
	c3bHtdRawData := bean.NeedClientSignTxData{}
	c3bHtdRawData.Hex = c3bHtdTx["hex"].(string)
	c3bHtdRawData.Inputs = c3bHtdTx["inputs"]
	c3bHtdRawData.IsMultisig = true
	c3bHtdRawData.PubKeyA = dataFromBob.PayeeCurrHtlcTempAddressPubKey
	c3bHtdRawData.PubKeyB = payerChannelPubKey
	needBobSignData.C3bHtlcHtdPartialData = c3bHtdRawData
	needAliceSignData.C3bHtlcHtdRawData = c3bHtdRawData
	//endregion

	//region 7 c3b htlc Hlock
	c3bHlockTx, err := createHtlcLockByHForBobAtPayerSide(*channelInfo, dataFromBob, c3bSignedHtlcHex, commitmentTransaction.HtlcH, payeeChannelPubKey, payerChannelPubKey, channelInfo.PropertyId, commitmentTransaction.AmountToHtlc)
	if err != nil {
		log.Println(err)
		return nil, errors.New("fail to create HlockHex for C3b")
	}
	needBobSignData.C3bHtlcHlockPartialData = c3bHlockTx
	needAliceSignData.C3bHtlcHlockRawData = c3bHlockTx
	//endregion

	//region 8 c3b htlc br
	c3bHtlcBrTx, err := rpcClient.OmniCreateRawTransactionUseUnsendInput(
		bobHtlcMultiAddress,
		bobHtlcOutputs,
		payerChannelAddress,
		channelInfo.FundingAddress,
		channelInfo.PropertyId,
		commitmentTransaction.AmountToHtlc,
		getBtcMinerAmount(channelInfo.BtcAmount),
		0,
		&bobHtlcRedeemScript)
	if err != nil {
		log.Println(err)
		return nil, errors.New("fail to create HTD1b for C3b")
	}
	c3bHtlcBrRawData := bean.NeedClientSignTxData{}
	c3bHtlcBrRawData.Hex = c3bHtlcBrTx["hex"].(string)
	c3bHtlcBrRawData.Inputs = c3bHtlcBrTx["inputs"]
	c3bHtlcBrRawData.IsMultisig = true
	c3bHtlcBrRawData.PubKeyA = dataFromBob.PayeeCurrHtlcTempAddressPubKey
	c3bHtlcBrRawData.PubKeyB = payerChannelPubKey
	needAliceSignData.C3bHtlcBrRawData = c3bHtlcBrRawData
	//endregion

	//region 9  c3a htlc ht htrd
	c3aHtHex := aliceSignedC3b.C3aHtlcHtCompleteSignedHex
	c3aHtMultiAddress, c3aHtRedeemScript, c3aHtAddrScriptPubKey, err := createMultiSig(dataFromBob.C3aHtlcTempAddressForHtPubKey, payeeChannelPubKey)
	if err != nil {
		return nil, err
	}
	c3aHtOutputs, err := getInputsForNextTxByParseTxHashVout(c3aHtHex, c3aHtMultiAddress, c3aHtAddrScriptPubKey, c3aHtRedeemScript)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	c3bHtRdTx, err := rpcClient.OmniCreateRawTransactionUseUnsendInput(
		c3aHtMultiAddress,
		c3aHtOutputs,
		payerChannelAddress,
		channelInfo.FundingAddress,
		channelInfo.PropertyId,
		commitmentTransaction.AmountToHtlc,
		getBtcMinerAmount(channelInfo.BtcAmount),
		1000,
		&c3aHtRedeemScript)
	if err != nil {
		log.Println(err)
		return nil, errors.New("fail to create HTD1b for C3b")
	}

	c3aHtRdRawData := bean.NeedClientSignTxData{}
	c3aHtRdRawData.Hex = c3bHtRdTx["hex"].(string)
	c3aHtRdRawData.Inputs = c3bHtRdTx["inputs"]
	c3aHtRdRawData.IsMultisig = true
	c3aHtRdRawData.PubKeyA = dataFromBob.C3aHtlcTempAddressForHtPubKey
	c3aHtRdRawData.PubKeyB = payeeChannelPubKey
	needBobSignData.C3aHtlcHtHex = c3aHtHex
	needBobSignData.C3aHtlcHtrdPartialData = c3aHtRdRawData
	needAliceSignData.C3aHtlcHtrdRawData = c3aHtRdRawData

	//endregion

	//region 10  c3a htlc ht htbr
	c3aHtBrTx, err := rpcClient.OmniCreateRawTransactionUseUnsendInput(
		c3aHtMultiAddress,
		c3aHtOutputs,
		payeeChannelAddress,
		channelInfo.FundingAddress,
		channelInfo.PropertyId,
		commitmentTransaction.AmountToHtlc,
		getBtcMinerAmount(channelInfo.BtcAmount),
		0,
		&c3aHtRedeemScript)
	if err != nil {
		log.Println(err)
		return nil, errors.New("fail to create HTD1b for C3b")
	}

	c3aHtBrRawData := bean.NeedClientSignTxData{}
	c3aHtBrRawData.Hex = c3aHtBrTx["hex"].(string)
	c3aHtBrRawData.Inputs = c3aHtBrTx["inputs"]
	c3aHtBrRawData.IsMultisig = true
	c3aHtBrRawData.PubKeyA = dataFromBob.C3aHtlcTempAddressForHtPubKey
	c3aHtBrRawData.PubKeyB = payeeChannelPubKey
	needBobSignData.C3aHtlcHtbrRawData = c3aHtBrRawData
	//endregion

	//region 11  c3a htlc hlock hed
	c3aHlockHex := aliceSignedC3b.C3aHtlcHlockCompleteSignedHex
	c3aHlockMultiAddress, c3aHlockRedeemScript, c3aHlockAddrScriptPubKey, err := createMultiSig(commitmentTransaction.HtlcH, payeeChannelPubKey)
	if err != nil {
		return nil, err
	}
	c3aHlocOutputs, err := getInputsForNextTxByParseTxHashVout(c3aHlockHex, c3aHlockMultiAddress, c3aHlockAddrScriptPubKey, c3aHlockRedeemScript)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	c3aHedTx, err := rpcClient.OmniCreateRawTransactionUseUnsendInput(
		c3aHlockMultiAddress,
		c3aHlocOutputs,
		payeeChannelAddress,
		payeeChannelAddress,
		channelInfo.PropertyId,
		commitmentTransaction.AmountToHtlc,
		getBtcMinerAmount(channelInfo.BtcAmount),
		0,
		&c3aHlockRedeemScript)
	if err != nil {
		log.Println(err)
		return nil, errors.New("fail to create Hed1b for C3a")
	}

	c3aHedRawData := bean.NeedClientSignTxData{}
	c3aHedRawData.Hex = c3aHedTx["hex"].(string)
	c3aHedRawData.Inputs = c3aHedTx["inputs"]
	c3aHedRawData.IsMultisig = true
	c3aHedRawData.PubKeyA = commitmentTransaction.HtlcH
	c3aHedRawData.PubKeyB = payeeChannelPubKey
	needBobSignData.C3aHtlcHedRawData = c3aHedRawData
	//endregion

	if service.tempDataSendTo42PAtAliceSide == nil {
		service.tempDataSendTo42PAtAliceSide = make(map[string]bean.NeedBobSignHtlcSubTxOfC3bP2p)
	}
	service.tempDataSendTo42PAtAliceSide[user.PeerId+"_"+channelId] = needBobSignData
	_ = tx.Commit()

	return needAliceSignData, nil
}

// step 8 alice 响应 100103号协议，更新alice的承诺交易，推送42号p2p协议
func (service *htlcForwardTxManager) OnAliceSignedC3bSubTxAtAliceSide(msg bean.RequestMessage, user bean.User) (interface{}, error) {
	aliceSignedC3b := bean.AliceSignHtlcSubTxOfC3bResult{}
	_ = json.Unmarshal([]byte(msg.Data), &aliceSignedC3b)

	channelId := aliceSignedC3b.ChannelId
	dataFromBob_41 := service.tempDataFrom41PAtAliceSide[user.PeerId+"_"+channelId]
	needBobSignData := service.tempDataSendTo42PAtAliceSide[user.PeerId+"_"+channelId]

	bobRsmcHexIsExist := false
	if len(needBobSignData.C3bRsmcRdPartialData.Hex) > 0 {
		if pass, _ := rpcClient.CheckMultiSign(false, aliceSignedC3b.C3bRsmcRdPartialSignedHex, 1); pass == false {
			return nil, errors.New(enum.Tips_common_wrong + "c3b_rsmc_rd_partial_signed_hex")
		}
		needBobSignData.C3bRsmcRdPartialData.Hex = aliceSignedC3b.C3bRsmcRdPartialSignedHex

		if pass, _ := rpcClient.CheckMultiSign(false, aliceSignedC3b.C3bRsmcBrPartialSignedHex, 1); pass == false {
			return nil, errors.New(enum.Tips_common_wrong + "c3b_rsmc_br_partial_signed_hex")
		}
		bobRsmcHexIsExist = true
	}

	if pass, _ := rpcClient.CheckMultiSign(false, aliceSignedC3b.C3bHtlcHtdPartialSignedHex, 1); pass == false {
		return nil, errors.New(enum.Tips_common_wrong + "c3b_htlc_htd_partial_signed_hex")
	}
	needBobSignData.C3bHtlcHtdPartialData.Hex = aliceSignedC3b.C3bHtlcHtdPartialSignedHex

	if pass, _ := rpcClient.CheckMultiSign(false, aliceSignedC3b.C3bHtlcBrPartialSignedHex, 1); pass == false {
		return nil, errors.New(enum.Tips_common_wrong + "c3b_htlc_br_partial_signed_hex")
	}

	if pass, _ := rpcClient.CheckMultiSign(false, aliceSignedC3b.C3bHtlcHlockPartialSignedHex, 1); pass == false {
		return nil, errors.New(enum.Tips_common_wrong + "c3b_htlc_hlock_partial_signed_hex")
	}
	needBobSignData.C3bHtlcHlockPartialData.Hex = aliceSignedC3b.C3bHtlcHlockPartialSignedHex

	if pass, _ := rpcClient.CheckMultiSign(false, aliceSignedC3b.C3aHtlcHtrdPartialSignedHex, 1); pass == false {
		return nil, errors.New(enum.Tips_common_wrong + "c3a_htlc_htrd_partial_signed_hex")
	}
	needBobSignData.C3aHtlcHtrdPartialData.Hex = aliceSignedC3b.C3aHtlcHtrdPartialSignedHex

	tx, err := user.Db.Begin(true)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	commitmentTransaction := &dao.CommitmentTransaction{}
	err = tx.Select(
		q.Eq("CurrHash", dataFromBob_41.PayerCommitmentTxHash),
		q.Eq("ChannelId", channelId)).First(commitmentTransaction)
	if err != nil {
		return nil, err
	}

	channelInfo := &dao.ChannelInfo{}
	err = tx.Select(
		q.Eq("ChannelId", channelId),
		q.Or(
			q.Eq("PeerIdA", user.PeerId),
			q.Eq("PeerIdB", user.PeerId))).
		First(channelInfo)
	if channelInfo == nil {
		return nil, errors.New("not found the channel " + channelId)
	}

	htlcRequestInfo := &dao.AddHtlcRequestInfo{}
	err = tx.Select(
		q.Eq("ChannelId", channelId),
		q.Eq("H", commitmentTransaction.HtlcH),
		q.Eq("PropertyId", channelInfo.PropertyId)).
		OrderBy("CreateAt").
		Reverse().
		First(htlcRequestInfo)
	if err != nil {
		return nil, err
	}

	htlcRequestInfo.CurrState = dao.NS_Finish
	_ = tx.Update(htlcRequestInfo)

	payerChannelPubKey := channelInfo.PubKeyA
	payerChannelAddress := channelInfo.AddressA
	if user.PeerId == channelInfo.PeerIdB {
		payerChannelPubKey = channelInfo.PubKeyB
		payerChannelAddress = channelInfo.AddressB
	}

	//region 处理付款方收到的已经签名C3a的子交易，及上一个BR的签名，RSMCBR，HBR的创建
	if commitmentTransaction.CurrState == dao.TxInfoState_Create {

		//region 1 根据对方传过来的上一个交易的临时rsmc私钥，签名最近的BR交易，保证对方确实放弃了上一个承诺交易
		err := signLastBR(tx, dao.BRType_Rmsc, *channelInfo, user.PeerId, dataFromBob_41.PayeeLastTempAddressPrivateKey, commitmentTransaction.LastCommitmentTxId)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		//endregion

		// region 2 保存c3b的rsmc的br到本地
		tempOtherSideCommitmentTx := &dao.CommitmentTransaction{}
		if bobRsmcHexIsExist {
			bobRsmcMultiAddress, bobRsmcRedeemScript, bobRsmcMultiAddressScriptPubKey, err := createMultiSig(dataFromBob_41.PayeeCurrRsmcTempAddressPubKey, payerChannelPubKey)
			if err != nil {
				return nil, err
			}
			c3bRsmcOutputs, err := getInputsForNextTxByParseTxHashVout(dataFromBob_41.C3bRsmcPartialSignedData.Hex, bobRsmcMultiAddress, bobRsmcMultiAddressScriptPubKey, bobRsmcRedeemScript)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			tempOtherSideCommitmentTx.Id = commitmentTransaction.Id
			tempOtherSideCommitmentTx.PropertyId = channelInfo.PropertyId
			tempOtherSideCommitmentTx.RSMCTempAddressPubKey = dataFromBob_41.PayeeCurrRsmcTempAddressPubKey
			tempOtherSideCommitmentTx.RSMCMultiAddress = bobRsmcMultiAddress
			tempOtherSideCommitmentTx.RSMCRedeemScript = bobRsmcRedeemScript
			tempOtherSideCommitmentTx.RSMCMultiAddressScriptPubKey = bobRsmcMultiAddressScriptPubKey
			tempOtherSideCommitmentTx.RSMCTxHex = dataFromBob_41.C3bRsmcPartialSignedData.Hex
			tempOtherSideCommitmentTx.RSMCTxid = rpcClient.GetTxId(tempOtherSideCommitmentTx.RSMCTxHex)
			tempOtherSideCommitmentTx.AmountToRSMC = commitmentTransaction.AmountToCounterparty
			err = createCurrCommitmentTxPartialSignedBR(tx, dao.BRType_Rmsc, channelInfo, tempOtherSideCommitmentTx, c3bRsmcOutputs, payerChannelAddress, aliceSignedC3b.C3bRsmcBrPartialSignedHex, user)
			if err != nil {
				log.Println(err)
				return nil, err
			}
		}
		//endregion

		// region 3 保存c3b的htlc的br到本地
		bobHtlcMultiAddress, bobHtlcRedeemScript, bobHtlcMultiAddressScriptPubKey, err := createMultiSig(dataFromBob_41.PayeeCurrHtlcTempAddressPubKey, payerChannelPubKey)
		if err != nil {
			return nil, err
		}

		c3bHtlcOutputs, err := getInputsForNextTxByParseTxHashVout(dataFromBob_41.C3bHtlcPartialSignedData.Hex, bobHtlcMultiAddress, bobHtlcMultiAddressScriptPubKey, bobHtlcRedeemScript)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		tempOtherSideCommitmentTx.Id = commitmentTransaction.Id
		tempOtherSideCommitmentTx.PropertyId = channelInfo.PropertyId
		tempOtherSideCommitmentTx.RSMCTempAddressPubKey = dataFromBob_41.PayeeCurrHtlcTempAddressPubKey
		tempOtherSideCommitmentTx.RSMCMultiAddress = bobHtlcMultiAddress
		tempOtherSideCommitmentTx.RSMCMultiAddressScriptPubKey = bobHtlcMultiAddressScriptPubKey
		tempOtherSideCommitmentTx.RSMCRedeemScript = bobHtlcRedeemScript
		tempOtherSideCommitmentTx.RSMCTxHex = dataFromBob_41.C3bHtlcPartialSignedData.Hex
		tempOtherSideCommitmentTx.RSMCTxid = rpcClient.GetTxId(tempOtherSideCommitmentTx.RSMCTxHex)
		tempOtherSideCommitmentTx.AmountToRSMC = commitmentTransaction.AmountToHtlc
		err = createCurrCommitmentTxPartialSignedBR(tx, dao.BRType_Htlc, channelInfo, tempOtherSideCommitmentTx, c3bHtlcOutputs, payerChannelAddress, aliceSignedC3b.C3bHtlcBrPartialSignedHex, user)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		//endregion

		//region 4 更新收到的签名交易
		_, _, err = checkHexAndUpdateC3aOn42Protocal(tx, dataFromBob_41, *htlcRequestInfo, *channelInfo, commitmentTransaction, user)
		if err != nil {
			return nil, err
		}
		//endregion
	}
	//endregion

	commitmentTransaction.RsmcRawTxData = bean.NeedClientSignTxData{}
	commitmentTransaction.HtlcRawTxData = bean.NeedClientSignTxData{}
	commitmentTransaction.ToCounterpartyRawTxData = bean.NeedClientSignTxData{}
	commitmentTransaction.CurrState = dao.TxInfoState_Htlc_WaitHTRD1aSign
	commitmentTransaction.SignAt = time.Now()
	_ = tx.Update(commitmentTransaction)

	_ = tx.Commit()

	return needBobSignData, nil
}

// step 9 bob 响应 42号协议 构造需要bob签名的数据，缓存来自42号协议的数据 推送110042
func (service *htlcForwardTxManager) OnGetNeedBobSignC3bSubTxAtBobSide(msgData string, user bean.User) (interface{}, error) {
	c3bCacheData := bean.NeedBobSignHtlcSubTxOfC3bP2p{}
	_ = json.Unmarshal([]byte(msgData), &c3bCacheData)

	if service.tempDataFrom42PAtBobSide == nil {
		service.tempDataFrom42PAtBobSide = make(map[string]bean.NeedBobSignHtlcSubTxOfC3bP2p)
	}
	service.tempDataFrom42PAtBobSide[user.PeerId+"_"+c3bCacheData.ChannelId] = c3bCacheData

	needBobSign := bean.NeedBobSignHtlcSubTxOfC3b{}
	needBobSign.ChannelId = c3bCacheData.ChannelId
	needBobSign.C3aHtlcHedRawData = c3bCacheData.C3aHtlcHedRawData
	needBobSign.C3aHtlcHtrdPartialData = c3bCacheData.C3aHtlcHtrdPartialData
	needBobSign.C3aHtlcHtbrRawData = c3bCacheData.C3aHtlcHtbrRawData
	needBobSign.C3bRsmcRdPartialData = c3bCacheData.C3bRsmcRdPartialData
	needBobSign.C3bHtlcHlockPartialData = c3bCacheData.C3bHtlcHlockPartialData
	needBobSign.C3bHtlcHtdPartialData = c3bCacheData.C3bHtlcHtdPartialData
	return needBobSign, nil
}

// step 10 bob 响应100104:缓存签名的结果，生成hlock的 he让bob继续签名
func (service *htlcForwardTxManager) OnBobSignedC3bSubTxAtBobSide(msg bean.RequestMessage, user bean.User) (interface{}, error) {

	jsonObj := bean.BobSignedHtlcSubTxOfC3b{}
	_ = json.Unmarshal([]byte(msg.Data), &jsonObj)

	if tool.CheckIsString(&jsonObj.ChannelId) == false {
		return nil, errors.New(enum.Tips_common_empty + "channel_id")
	}

	if tool.CheckIsString(&jsonObj.CurrHtlcTempAddressForHePubKey) == false {
		return nil, errors.New(enum.Tips_common_empty + "curr_htlc_temp_address_for_he_pub_key")
	}

	c3bCacheData := service.tempDataFrom42PAtBobSide[user.PeerId+"_"+jsonObj.ChannelId]

	if pass, _ := rpcClient.CheckMultiSign(false, jsonObj.C3aHtlcHtrdCompleteSignedHex, 2); pass == false {
		return nil, errors.New("error sign c3a_htlc_htrd_complete_signed_hex")
	}
	c3bCacheData.C3aHtlcHtrdPartialData.Hex = jsonObj.C3aHtlcHtrdCompleteSignedHex

	if pass, _ := rpcClient.CheckMultiSign(false, jsonObj.C3aHtlcHedPartialSignedHex, 1); pass == false {
		return nil, errors.New("error sign c3a_htlc_hed_partial_signed_hex")
	}
	c3bCacheData.C3aHtlcHedRawData.Hex = jsonObj.C3aHtlcHedPartialSignedHex

	if pass, _ := rpcClient.CheckMultiSign(false, jsonObj.C3aHtlcHtbrPartialSignedHex, 1); pass == false {
		return nil, errors.New("error sign c3a_htlc_htbr_partial_signed_hex")
	}
	c3bCacheData.C3aHtlcHtbrRawData.Hex = jsonObj.C3aHtlcHtbrPartialSignedHex

	if pass, _ := rpcClient.CheckMultiSign(false, jsonObj.C3bRsmcRdCompleteSignedHex, 2); pass == false {
		return nil, errors.New("error sign c3b_rsmc_rd_complete_signed_hex")
	}
	c3bCacheData.C3bRsmcRdPartialData.Hex = jsonObj.C3bRsmcRdCompleteSignedHex

	if pass, _ := rpcClient.CheckMultiSign(false, jsonObj.C3bHtlcHlockCompleteSignedHex, 2); pass == false {
		return nil, errors.New("error sign c3b_htlc_hlock_complete_signed_hex")
	}
	c3bCacheData.C3bHtlcHlockPartialData.Hex = jsonObj.C3bHtlcHlockCompleteSignedHex

	if pass, _ := rpcClient.CheckMultiSign(false, jsonObj.C3bHtlcHtdCompleteSignedHex, 2); pass == false {
		return nil, errors.New("error sign c3b_htlc_htd_complete_signed_hex")
	}
	c3bCacheData.C3bHtlcHtdPartialData.Hex = jsonObj.C3bHtlcHtdCompleteSignedHex

	service.tempDataFrom42PAtBobSide[user.PeerId+"_"+jsonObj.ChannelId] = c3bCacheData

	//根据Hlock，创建C3b的He H+bobChannelPubkey作为输入，aliceChannelPubkey+bob的临时地址3作为输出
	tx, err := user.Db.Begin(true)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer tx.Rollback()

	latestCommitmentTx, err := getLatestCommitmentTxUseDbTx(tx, jsonObj.ChannelId, user.PeerId)
	if err != nil {
		return nil, err
	}

	channelInfo := &dao.ChannelInfo{}
	err = tx.Select(
		q.Eq("ChannelId", latestCommitmentTx.ChannelId),
		q.Or(
			q.Eq("PeerIdA", user.PeerId),
			q.Eq("PeerIdB", user.PeerId))).
		First(channelInfo)
	if channelInfo == nil {
		return nil, errors.New("not found the channel " + jsonObj.ChannelId)
	}

	c3bHeRawData, err := saveHtlcHeTxForPayee(tx, *channelInfo, latestCommitmentTx, jsonObj, jsonObj.C3bHtlcHlockCompleteSignedHex, user)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	needBobSignData := bean.NeedBobSignHtlcHeTxOfC3b{}
	needBobSignData.ChannelId = jsonObj.ChannelId
	needBobSignData.C3bHtlcHlockHeRawData = *c3bHeRawData
	needBobSignData.PayerPeerId = c3bCacheData.PayerPeerId
	needBobSignData.PayerNodeAddress = c3bCacheData.PayerNodeAddress
	return needBobSignData, nil
}

// step 11 bob 响应100105:收款方完成Hlock的he的部分签名，更新C3b的信息，最后推送43给alice和推送正向H的创建htlc的结果
func (service *htlcForwardTxManager) OnBobSignHtRdAtBobSide_42(msgData string, user bean.User) (toAlice, toBob interface{}, err error) {
	jsonObj := bean.BobSignedHtlcHeTxOfC3b{}
	_ = json.Unmarshal([]byte(msgData), &jsonObj)

	if tool.CheckIsString(&jsonObj.ChannelId) == false {
		return nil, nil, errors.New(enum.Tips_common_empty + "channel_id")
	}

	c3bCacheData := service.tempDataFrom42PAtBobSide[user.PeerId+"_"+jsonObj.ChannelId]
	if len(c3bCacheData.ChannelId) == 0 {
		return nil, nil, errors.New(enum.Tips_common_wrong + "channel_id")
	}

	if pass, _ := rpcClient.CheckMultiSign(false, jsonObj.C3bHtlcHlockHePartialSignedHex, 1); pass == false {
		return nil, nil, errors.New("error sign c3b_htlc_hlock_he_partial_signed_hex")
	}

	tx, err := user.Db.Begin(true)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	defer tx.Rollback()

	c3aRetData := bean.C3aSignedHerdTxOfC3bP2p{}
	c3aRetData.ChannelId = c3bCacheData.ChannelId

	latestCommitmentTx := &dao.CommitmentTransaction{}
	err = tx.Select(q.Eq("CurrHash", c3bCacheData.PayeeCommitmentTxHash)).First(latestCommitmentTx)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	channelInfo := &dao.ChannelInfo{}
	err = tx.Select(q.Eq("ChannelId", latestCommitmentTx.ChannelId)).First(channelInfo)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	bobChannelPubKey := channelInfo.PubKeyB
	bobChannelAddress := channelInfo.AddressB
	if user.PeerId == channelInfo.PeerIdA {
		bobChannelPubKey = channelInfo.PubKeyA
		bobChannelAddress = channelInfo.AddressA
	}

	//region 1 返回给alice的htrd签名数据
	c3aRetData.C3aHtlcHtrdCompleteSignedHex = c3bCacheData.C3aHtlcHtrdPartialData.Hex
	//endregion
	c3aRetData.C3aHtlcHedPartialSignedHex = c3bCacheData.C3aHtlcHedRawData.Hex

	if latestCommitmentTx.CurrState == dao.TxInfoState_Create {
		//创建ht的BR
		c3aHtMultiAddress, c3aHtRedeemScript, c3aHtMultiAddressScriptPubKey, err := createMultiSig(c3bCacheData.C3aHtlcTempAddressForHtPubKey, bobChannelPubKey)
		if err != nil {
			return nil, nil, err
		}

		c3aHtOutputs, err := getInputsForNextTxByParseTxHashVout(c3bCacheData.C3aHtlcHtHex, c3aHtMultiAddress, c3aHtMultiAddressScriptPubKey, c3aHtRedeemScript)
		if err != nil {
			return nil, nil, err
		}

		tempOtherSideCommitmentTx := &dao.CommitmentTransaction{}
		tempOtherSideCommitmentTx.Id = latestCommitmentTx.Id
		tempOtherSideCommitmentTx.PropertyId = channelInfo.PropertyId
		tempOtherSideCommitmentTx.RSMCTempAddressPubKey = c3bCacheData.C3aHtlcTempAddressForHtPubKey
		tempOtherSideCommitmentTx.RSMCMultiAddress = c3aHtMultiAddress
		tempOtherSideCommitmentTx.RSMCRedeemScript = c3aHtRedeemScript
		tempOtherSideCommitmentTx.RSMCMultiAddressScriptPubKey = c3aHtMultiAddressScriptPubKey
		tempOtherSideCommitmentTx.RSMCTxHex = c3bCacheData.C3aHtlcHtHex
		tempOtherSideCommitmentTx.RSMCTxid = c3aHtOutputs[0].Txid
		tempOtherSideCommitmentTx.AmountToRSMC = latestCommitmentTx.AmountToHtlc
		err = createCurrCommitmentTxPartialSignedBR(tx, dao.BRType_Ht1a, channelInfo, tempOtherSideCommitmentTx, c3aHtOutputs, bobChannelAddress, c3bCacheData.C3aHtlcHtbrRawData.Hex, user)
		if err != nil {
			log.Println(err)
			return nil, nil, err
		}

		latestCommitmentTx, _, err = checkHexAndUpdateC3bOn42Protocal(tx, c3bCacheData, *channelInfo, latestCommitmentTx, user)
		if err != nil {
			log.Println(err.Error())
			return nil, nil, err
		}

		//更新He交易
		_ = updateHtlcHeTxForPayee(tx, *channelInfo, latestCommitmentTx, jsonObj.C3bHtlcHlockHePartialSignedHex)
	}

	channelInfo.CurrState = dao.ChannelState_HtlcTx
	_ = tx.Update(channelInfo)

	_ = tx.Commit()

	key := user.PeerId + "_" + channelInfo.ChannelId
	delete(service.tempDataFrom42PAtBobSide, key)
	return c3aRetData, latestCommitmentTx, nil
}

// step 12 响应43号协议: 保存htrd和hed 推送110043给Alice
func (service *htlcForwardTxManager) OnGetHtrdTxDataFromBobAtAliceSide_43(msgData string, user bean.User) (data interface{}, err error) {

	c3aHtrdData := bean.C3aSignedHerdTxOfC3bP2p{}
	_ = json.Unmarshal([]byte(msgData), &c3aHtrdData)

	if tool.CheckIsString(&c3aHtrdData.ChannelId) == false {
		return nil, errors.New(enum.Tips_common_empty + "channel_id")
	}

	if pass, _ := rpcClient.CheckMultiSign(false, c3aHtrdData.C3aHtlcHedPartialSignedHex, 1); pass == false {
		return nil, errors.New("error sign c3a_htlc_hed_partial_data")
	}

	if pass, _ := rpcClient.CheckMultiSign(false, c3aHtrdData.C3aHtlcHtrdCompleteSignedHex, 2); pass == false {
		return nil, errors.New("error sign c3a_htlc_htrd_complete_signed_hex")
	}

	tx, err := user.Db.Begin(true)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer tx.Rollback()

	latestCommitmentTx, err := getLatestCommitmentTxUseDbTx(tx, c3aHtrdData.ChannelId, user.PeerId)
	if err != nil {
		return nil, err
	}

	channelInfo := dao.ChannelInfo{}
	err = tx.Select(
		q.Eq("ChannelId", latestCommitmentTx.ChannelId),
		q.Or(
			q.Eq("PeerIdA", user.PeerId),
			q.Eq("PeerIdB", user.PeerId))).
		First(&channelInfo)
	if &channelInfo == nil {
		return nil, errors.New("not found the channel " + c3aHtrdData.ChannelId)
	}

	if latestCommitmentTx.CurrState == dao.TxInfoState_Htlc_WaitHTRD1aSign {
		_, err := saveHtRD1a(tx, c3aHtrdData.C3aHtlcHtrdCompleteSignedHex, *latestCommitmentTx, user)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		_ = createHed1a(tx, c3aHtrdData.C3aHtlcHedPartialSignedHex, channelInfo, *latestCommitmentTx, user)
	}

	latestCommitmentTx.RsmcRawTxData = bean.NeedClientSignTxData{}
	latestCommitmentTx.HtlcRawTxData = bean.NeedClientSignTxData{}
	latestCommitmentTx.ToCounterpartyRawTxData = bean.NeedClientSignTxData{}
	latestCommitmentTx.CurrState = dao.TxInfoState_Htlc_GetH
	_ = tx.Update(latestCommitmentTx)

	channelInfo.CurrState = dao.ChannelState_HtlcTx
	_ = tx.Update(channelInfo)
	tx.Commit()

	//同步通道信息到tracker
	sendChannelStateToTracker(channelInfo, *latestCommitmentTx)

	key := user.PeerId + "_" + channelInfo.ChannelId
	delete(service.tempDataSendTo40PAtAliceSide, key)
	delete(service.tempDataFrom41PAtAliceSide, key)
	delete(service.tempDataSendTo42PAtAliceSide, key)

	return latestCommitmentTx, nil
}

// 创建付款方C3a
func htlcPayerCreateCommitmentTx_C3a(tx storm.Node, channelInfo *dao.ChannelInfo, requestData bean.CreateHtlcTxForC3a, totalStep int, currStep int, latestCommitmentTx *dao.CommitmentTransaction, user bean.User) (*dao.CommitmentTransaction, error) {

	fundingTransaction := getFundingTransactionByChannelId(tx, channelInfo.ChannelId, user.PeerId)
	if fundingTransaction == nil {
		return nil, errors.New("not found fundingTransaction")
	}
	// htlc的资产分配方案
	var outputBean = commitmentTxOutputBean{}
	amountAndFee, _ := decimal.NewFromFloat(requestData.Amount).Mul(decimal.NewFromFloat(1 + config.GetHtlcFee()*float64(totalStep-(currStep+1)))).Round(8).Float64()
	outputBean.RsmcTempPubKey = requestData.CurrRsmcTempAddressPubKey
	outputBean.HtlcTempPubKey = requestData.CurrHtlcTempAddressPubKey

	aliceIsPayer := true
	if user.PeerId == channelInfo.PeerIdB {
		aliceIsPayer = false
	}
	outputBean.AmountToHtlc = amountAndFee
	if aliceIsPayer { //Alice pay money to bob Alice是付款方
		outputBean.AmountToRsmc, _ = decimal.NewFromFloat(fundingTransaction.AmountA).Sub(decimal.NewFromFloat(amountAndFee)).Round(8).Float64()
		outputBean.AmountToCounterparty = fundingTransaction.AmountB
		outputBean.OppositeSideChannelPubKey = channelInfo.PubKeyB
		outputBean.OppositeSideChannelAddress = channelInfo.AddressB
	} else { //	bob pay money to alice
		outputBean.AmountToRsmc, _ = decimal.NewFromFloat(fundingTransaction.AmountB).Sub(decimal.NewFromFloat(amountAndFee)).Round(8).Float64()
		outputBean.AmountToCounterparty = fundingTransaction.AmountA
		outputBean.OppositeSideChannelPubKey = channelInfo.PubKeyA
		outputBean.OppositeSideChannelAddress = channelInfo.AddressA
	}
	if latestCommitmentTx.Id > 0 {
		outputBean.AmountToRsmc, _ = decimal.NewFromFloat(latestCommitmentTx.AmountToRSMC).Sub(decimal.NewFromFloat(amountAndFee)).Round(8).Float64()
		outputBean.AmountToCounterparty = latestCommitmentTx.AmountToCounterparty
	}

	newCommitmentTxInfo, err := createCommitmentTx(user.PeerId, channelInfo, fundingTransaction, outputBean, &user)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	newCommitmentTxInfo.TxType = dao.CommitmentTransactionType_Htlc
	newCommitmentTxInfo.RSMCTempAddressIndex = requestData.CurrRsmcTempAddressIndex
	newCommitmentTxInfo.HTLCTempAddressIndex = requestData.CurrHtlcTempAddressIndex

	allUsedTxidTemp := ""
	// rsmc
	if newCommitmentTxInfo.AmountToRSMC > 0 {
		rsmcTxData, usedTxid, err := rpcClient.OmniCreateRawTransactionUseSingleInput(
			int(newCommitmentTxInfo.TxType),
			channelInfo.ChannelAddress,
			newCommitmentTxInfo.RSMCMultiAddress,
			channelInfo.PropertyId,
			newCommitmentTxInfo.AmountToRSMC,
			0,
			0, &channelInfo.ChannelAddressRedeemScript, "")
		if err != nil {
			log.Println(err)
			return nil, err
		}
		allUsedTxidTemp += usedTxid
		newCommitmentTxInfo.RsmcInputTxid = usedTxid
		newCommitmentTxInfo.RSMCTxHex = rsmcTxData["hex"].(string)

		signHexData := bean.NeedClientSignTxData{}
		signHexData.Hex = newCommitmentTxInfo.RSMCTxHex
		signHexData.Inputs = rsmcTxData["inputs"]
		signHexData.IsMultisig = true
		signHexData.PubKeyA = channelInfo.PubKeyA
		signHexData.PubKeyB = channelInfo.PubKeyB
		newCommitmentTxInfo.RsmcRawTxData = signHexData
	}

	//htlc
	if newCommitmentTxInfo.AmountToHtlc > 0 {
		htlcTxData, usedTxid, err := rpcClient.OmniCreateRawTransactionUseSingleInput(
			int(newCommitmentTxInfo.TxType),
			channelInfo.ChannelAddress,
			newCommitmentTxInfo.HTLCMultiAddress,
			channelInfo.PropertyId,
			newCommitmentTxInfo.AmountToHtlc,
			0,
			0, &channelInfo.ChannelAddressRedeemScript, allUsedTxidTemp)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		allUsedTxidTemp += "," + usedTxid
		newCommitmentTxInfo.HtlcRoutingPacket = requestData.RoutingPacket

		currBlockHeight, err := rpcClient.GetBlockCount()
		if err != nil {
			return nil, errors.New("fail to get blockHeight ,please try again later")
		}
		newCommitmentTxInfo.HtlcCltvExpiry = requestData.CltvExpiry
		newCommitmentTxInfo.BeginBlockHeight = currBlockHeight

		newCommitmentTxInfo.HtlcTxHex = htlcTxData["hex"].(string)
		newCommitmentTxInfo.HtlcH = requestData.H
		if aliceIsPayer {
			newCommitmentTxInfo.HtlcSender = channelInfo.PeerIdA
		} else {
			newCommitmentTxInfo.HtlcSender = channelInfo.PeerIdB
		}

		signHexData := bean.NeedClientSignTxData{}
		signHexData.Hex = newCommitmentTxInfo.HtlcTxHex
		signHexData.Inputs = htlcTxData["inputs"]
		signHexData.IsMultisig = true
		signHexData.PubKeyA = channelInfo.PubKeyA
		signHexData.PubKeyB = channelInfo.PubKeyB
		newCommitmentTxInfo.HtlcRawTxData = signHexData

	}

	//create to Bob tx
	if newCommitmentTxInfo.AmountToCounterparty > 0 {
		toBobTxData, err := rpcClient.OmniCreateRawTransactionUseRestInput(
			int(newCommitmentTxInfo.TxType),
			channelInfo.ChannelAddress,
			allUsedTxidTemp,
			outputBean.OppositeSideChannelAddress,
			channelInfo.FundingAddress,
			channelInfo.PropertyId,
			newCommitmentTxInfo.AmountToCounterparty,
			getBtcMinerAmount(channelInfo.BtcAmount),
			&channelInfo.ChannelAddressRedeemScript)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		newCommitmentTxInfo.ToCounterpartyTxHex = toBobTxData["hex"].(string)

		signHexData := bean.NeedClientSignTxData{}
		signHexData.Hex = newCommitmentTxInfo.ToCounterpartyTxHex
		signHexData.Inputs = toBobTxData["inputs"]
		signHexData.IsMultisig = true
		signHexData.PubKeyA = channelInfo.PubKeyA
		signHexData.PubKeyB = channelInfo.PubKeyB
		newCommitmentTxInfo.ToCounterpartyRawTxData = signHexData

	}

	newCommitmentTxInfo.CurrState = dao.TxInfoState_Init
	newCommitmentTxInfo.LastHash = ""
	newCommitmentTxInfo.CurrHash = ""
	if latestCommitmentTx.Id > 0 {
		newCommitmentTxInfo.LastCommitmentTxId = latestCommitmentTx.Id
		newCommitmentTxInfo.LastHash = latestCommitmentTx.CurrHash
	}
	err = tx.Save(newCommitmentTxInfo)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	bytes, err := json.Marshal(newCommitmentTxInfo)
	msgHash := tool.SignMsgWithSha256(bytes)
	newCommitmentTxInfo.CurrHash = msgHash
	err = tx.Update(newCommitmentTxInfo)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return newCommitmentTxInfo, nil
}

// 创建收款方C3b
func htlcPayeeCreateCommitmentTx_C3b(tx storm.Node, channelInfo *dao.ChannelInfo, reqData bean.BobSignedC3a, payerData bean.CreateHtlcTxForC3aOfP2p, latestCommitmentTx *dao.CommitmentTransaction, signedToOtherHex string, user bean.User) (*dao.CommitmentTransaction, error) {

	channelIds := strings.Split(payerData.RoutingPacket, ",")
	var totalStep = len(channelIds)
	var currStep = 0
	for index, channelId := range channelIds {
		if channelId == channelInfo.ChannelId {
			currStep = index
			break
		}
	}
	fundingTransaction := getFundingTransactionByChannelId(tx, channelInfo.ChannelId, user.PeerId)
	if fundingTransaction == nil {
		return nil, errors.New("not found fundingTransaction")
	}

	// htlc的资产分配方案
	var outputBean = commitmentTxOutputBean{}
	decimal.DivisionPrecision = 8
	amountAndFee, _ := decimal.NewFromFloat(payerData.Amount).Mul(decimal.NewFromFloat((1 + config.GetHtlcFee()*float64(totalStep-(currStep+1))))).Round(8).Float64()
	outputBean.RsmcTempPubKey = reqData.CurrRsmcTempAddressPubKey
	outputBean.HtlcTempPubKey = reqData.CurrHtlcTempAddressPubKey

	bobIsPayee := true
	if user.PeerId == channelInfo.PeerIdA {
		bobIsPayee = false
	}
	outputBean.AmountToHtlc = amountAndFee
	if bobIsPayee { //Alice pay money to bob
		outputBean.AmountToRsmc = fundingTransaction.AmountB
		outputBean.AmountToCounterparty, _ = decimal.NewFromFloat(fundingTransaction.AmountA).Sub(decimal.NewFromFloat(amountAndFee)).Round(8).Float64()
		outputBean.OppositeSideChannelPubKey = channelInfo.PubKeyA
		outputBean.OppositeSideChannelAddress = channelInfo.AddressA
	} else { //	bob pay money to alice
		outputBean.AmountToRsmc = fundingTransaction.AmountA
		outputBean.AmountToCounterparty, _ = decimal.NewFromFloat(fundingTransaction.AmountB).Sub(decimal.NewFromFloat(amountAndFee)).Round(8).Float64()
		outputBean.OppositeSideChannelPubKey = channelInfo.PubKeyB
		outputBean.OppositeSideChannelAddress = channelInfo.AddressB
	}
	if latestCommitmentTx.Id > 0 {
		outputBean.AmountToCounterparty, _ = decimal.NewFromFloat(latestCommitmentTx.AmountToCounterparty).Sub(decimal.NewFromFloat(amountAndFee)).Round(8).Float64()
		outputBean.AmountToRsmc = latestCommitmentTx.AmountToRSMC
	}

	newCommitmentTxInfo, err := createCommitmentTx(user.PeerId, channelInfo, fundingTransaction, outputBean, &user)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	newCommitmentTxInfo.FromCounterpartySideForMeTxHex = signedToOtherHex
	newCommitmentTxInfo.TxType = dao.CommitmentTransactionType_Htlc
	newCommitmentTxInfo.RSMCTempAddressIndex = reqData.CurrRsmcTempAddressIndex
	newCommitmentTxInfo.HTLCTempAddressIndex = reqData.CurrHtlcTempAddressIndex

	allUsedTxidTemp := ""
	// rsmc
	if newCommitmentTxInfo.AmountToRSMC > 0 {
		rsmcTxData, usedTxid, err := rpcClient.OmniCreateRawTransactionUseSingleInput(
			int(newCommitmentTxInfo.TxType),
			channelInfo.ChannelAddress,
			newCommitmentTxInfo.RSMCMultiAddress,
			channelInfo.PropertyId,
			newCommitmentTxInfo.AmountToRSMC,
			0,
			0, &channelInfo.ChannelAddressRedeemScript, "")
		if err != nil {
			log.Println(err)
			return nil, err
		}
		allUsedTxidTemp += usedTxid
		newCommitmentTxInfo.RsmcInputTxid = usedTxid
		newCommitmentTxInfo.RSMCTxHex = rsmcTxData["hex"].(string)

		signHexData := bean.NeedClientSignTxData{}
		signHexData.Hex = newCommitmentTxInfo.RSMCTxHex
		signHexData.Inputs = rsmcTxData["inputs"]
		signHexData.IsMultisig = true
		signHexData.PubKeyA = channelInfo.PubKeyA
		signHexData.PubKeyB = channelInfo.PubKeyB
		newCommitmentTxInfo.RsmcRawTxData = signHexData
	}

	// htlc
	if newCommitmentTxInfo.AmountToHtlc > 0 {
		htlcTxData, usedTxid, err := rpcClient.OmniCreateRawTransactionUseSingleInput(
			int(newCommitmentTxInfo.TxType),
			channelInfo.ChannelAddress,
			newCommitmentTxInfo.HTLCMultiAddress,
			channelInfo.PropertyId,
			newCommitmentTxInfo.AmountToHtlc,
			0,
			0, &channelInfo.ChannelAddressRedeemScript, allUsedTxidTemp)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		allUsedTxidTemp += "," + usedTxid
		newCommitmentTxInfo.HtlcRoutingPacket = payerData.RoutingPacket
		currBlockHeight, err := rpcClient.GetBlockCount()
		if err != nil {
			return nil, errors.New("fail to get blockHeight ,please try again later")
		}
		newCommitmentTxInfo.HtlcCltvExpiry = payerData.CltvExpiry
		newCommitmentTxInfo.BeginBlockHeight = currBlockHeight
		newCommitmentTxInfo.HtlcTxHex = htlcTxData["hex"].(string)

		signHexData := bean.NeedClientSignTxData{}
		signHexData.Hex = newCommitmentTxInfo.HtlcTxHex
		signHexData.Inputs = htlcTxData["inputs"]
		signHexData.IsMultisig = true
		signHexData.PubKeyA = channelInfo.PubKeyA
		signHexData.PubKeyB = channelInfo.PubKeyB
		newCommitmentTxInfo.HtlcRawTxData = signHexData

		newCommitmentTxInfo.HtlcH = payerData.H
		if bobIsPayee {
			newCommitmentTxInfo.HtlcSender = channelInfo.PeerIdA
		} else {
			newCommitmentTxInfo.HtlcSender = channelInfo.PeerIdB
		}
	}

	//create for other side tx
	if newCommitmentTxInfo.AmountToCounterparty > 0 {
		toBobTxData, err := rpcClient.OmniCreateRawTransactionUseRestInput(
			int(newCommitmentTxInfo.TxType),
			channelInfo.ChannelAddress,
			allUsedTxidTemp,
			outputBean.OppositeSideChannelAddress,
			channelInfo.FundingAddress,
			channelInfo.PropertyId,
			newCommitmentTxInfo.AmountToCounterparty,
			getBtcMinerAmount(channelInfo.BtcAmount),
			&channelInfo.ChannelAddressRedeemScript)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		newCommitmentTxInfo.ToCounterpartyTxHex = toBobTxData["hex"].(string)

		signHexData := bean.NeedClientSignTxData{}
		signHexData.Hex = newCommitmentTxInfo.ToCounterpartyTxHex
		signHexData.Inputs = toBobTxData["inputs"]
		signHexData.IsMultisig = true
		signHexData.PubKeyA = channelInfo.PubKeyA
		signHexData.PubKeyB = channelInfo.PubKeyB
		newCommitmentTxInfo.ToCounterpartyRawTxData = signHexData
	}

	newCommitmentTxInfo.CurrState = dao.TxInfoState_Init
	newCommitmentTxInfo.LastHash = ""
	newCommitmentTxInfo.CurrHash = ""
	if latestCommitmentTx.Id > 0 {
		newCommitmentTxInfo.LastCommitmentTxId = latestCommitmentTx.Id
		newCommitmentTxInfo.LastHash = latestCommitmentTx.CurrHash
	}
	err = tx.Save(newCommitmentTxInfo)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	bytes, err := json.Marshal(newCommitmentTxInfo)
	msgHash := tool.SignMsgWithSha256(bytes)
	newCommitmentTxInfo.CurrHash = msgHash
	err = tx.Update(newCommitmentTxInfo)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return newCommitmentTxInfo, nil
}

// 付款方更新C3a的信息
func checkHexAndUpdateC3aOn42Protocal(tx storm.Node, jsonObj bean.NeedAliceSignHtlcTxOfC3bP2p, htlcRequestInfo dao.AddHtlcRequestInfo, channelInfo dao.ChannelInfo, commitmentTransaction *dao.CommitmentTransaction, user bean.User) (retData *string, needNoticePayee bool, err error) {

	payeePubKey := channelInfo.PubKeyB
	payerAddress := channelInfo.AddressA
	otherSideAddress := channelInfo.AddressB
	if user.PeerId == channelInfo.PeerIdB {
		payerAddress = channelInfo.AddressB
		otherSideAddress = channelInfo.AddressA
		payeePubKey = channelInfo.PubKeyA
	}

	//region 1、检测 signedToOtherHex
	if len(commitmentTransaction.ToCounterpartyTxHex) > 0 {
		signedToCounterpartyHex := jsonObj.C3aCompleteSignedCounterpartyHex
		if tool.CheckIsString(&signedToCounterpartyHex) == false {
			err = errors.New("signedToOtherHex is empty at 41 protocol")
			log.Println(err)
			return nil, true, err
		}

		result, err := rpcClient.OmniDecodeTransaction(signedToCounterpartyHex)
		if err != nil {
			return nil, true, err
		}
		hexJsonObj := gjson.Parse(result)
		if channelInfo.ChannelAddress != hexJsonObj.Get("sendingaddress").String() {
			err = errors.New("wrong inputAddress at signedToOtherHex  at 41 protocol")
			log.Println(err)
			return nil, true, err
		}
		if channelInfo.PropertyId != hexJsonObj.Get("propertyid").Int() {
			err = errors.New("wrong propertyId at signedToOtherHex  at 41 protocol")
			log.Println(err)
			return nil, true, err
		}

		if otherSideAddress != hexJsonObj.Get("referenceaddress").String() {
			err = errors.New("wrong outputAddress at signedToOtherHex  at 41 protocol")
			log.Println(err)
			return nil, true, err
		}
		if commitmentTransaction.AmountToCounterparty != hexJsonObj.Get("amount").Float() {
			err = errors.New("wrong amount at signedToOtherHex  at 41 protocol")
			log.Println(err)
			return nil, true, err
		}
		commitmentTransaction.ToCounterpartyTxHex = signedToCounterpartyHex
		commitmentTransaction.ToCounterpartyTxid = hexJsonObj.Get("txid").String()
	}
	//endregion

	//region 2、检测 signedRsmcHex
	signedRsmcHex := jsonObj.C3aCompleteSignedRsmcHex
	if tool.CheckIsString(&signedRsmcHex) {
		result, err := rpcClient.OmniDecodeTransaction(signedRsmcHex)
		if err != nil {
			return nil, true, err
		}
		hexJsonObj := gjson.Parse(result)

		if channelInfo.ChannelAddress != hexJsonObj.Get("sendingaddress").String() {
			err = errors.New("wrong inputAddress at signedRsmcHex  at 41 protocol")
			log.Println(err)
			return nil, true, err
		}
		if channelInfo.PropertyId != hexJsonObj.Get("propertyid").Int() {
			err = errors.New("wrong propertyId at signedRsmcHex  at 41 protocol")
			log.Println(err)
			return nil, true, err
		}

		if commitmentTransaction.RSMCMultiAddress != hexJsonObj.Get("referenceaddress").String() {
			err = errors.New("wrong outputAddress at signedRsmcHex  at 41 protocol")
			log.Println(err)
			return nil, true, err
		}
		if commitmentTransaction.AmountToRSMC != hexJsonObj.Get("amount").Float() {
			err = errors.New("wrong amount at signedRsmcHex  at 41 protocol")
			log.Println(err)
			return nil, true, err
		}
		commitmentTransaction.RSMCTxHex = signedRsmcHex
		commitmentTransaction.RSMCTxid = hexJsonObj.Get("txid").String()
	}
	//endregion

	//region 3、检测 signedHtlcHex
	signedHtlcHex := jsonObj.C3aCompleteSignedHtlcHex
	if tool.CheckIsString(&signedHtlcHex) == false {
		err = errors.New("signedHtlcHex is empty at 41 protocol")
		log.Println(err)
		return nil, true, err
	}

	result, err := rpcClient.OmniDecodeTransaction(signedHtlcHex)
	if err != nil {
		return nil, true, err
	}
	hexJsonObj := gjson.Parse(result)
	if channelInfo.ChannelAddress != hexJsonObj.Get("sendingaddress").String() {
		err = errors.New("wrong inputAddress at signedHtlcHex  at 41 protocol")
		log.Println(err)
		return nil, true, err
	}
	if channelInfo.PropertyId != hexJsonObj.Get("propertyid").Int() {
		err = errors.New("wrong propertyId at signedHtlcHex  at 41 protocol")
		log.Println(err)
		return nil, true, err
	}

	if commitmentTransaction.HTLCMultiAddress != hexJsonObj.Get("referenceaddress").String() {
		err = errors.New("wrong outputAddress at signedHtlcHex  at 41 protocol")
		log.Println(err)
		return nil, true, err
	}
	if commitmentTransaction.AmountToHtlc != hexJsonObj.Get("amount").Float() {
		err = errors.New("wrong amount at signedHtlcHex  at 41 protocol")
		log.Println(err)
		return nil, true, err
	}
	commitmentTransaction.HtlcTxHex = signedHtlcHex
	commitmentTransaction.HTLCTxid = hexJsonObj.Get("txid").String()
	//endregion

	//region 4、rsmc Rd的保存
	payerRsmcRdHex := jsonObj.C3aRsmcRdPartialSignedData.Hex
	if tool.CheckIsString(&payerRsmcRdHex) {
		payerRDInputsFromRsmc, err := getInputsForNextTxByParseTxHashVout(
			signedRsmcHex,
			commitmentTransaction.RSMCMultiAddress,
			commitmentTransaction.RSMCMultiAddressScriptPubKey,
			commitmentTransaction.RSMCRedeemScript)
		if err != nil {
			log.Println(err)
			return nil, true, err
		}
		result, err = rpcClient.OmniDecodeTransactionWithPrevTxs(payerRsmcRdHex, payerRDInputsFromRsmc)
		if err != nil {
			log.Println(err)
			return nil, true, err
		}
		hexJsonObj = gjson.Parse(result)
		if commitmentTransaction.RSMCMultiAddress != hexJsonObj.Get("sendingaddress").String() {
			err = errors.New("wrong inputAddress at payerRsmcRdHex  at 41 protocol")
			log.Println(err)
			return nil, true, err
		}

		if payerAddress != hexJsonObj.Get("referenceaddress").String() {
			err = errors.New("wrong outputAddress at payerRsmcRdHex  at 41 protocol")
			log.Println(err)
			return nil, true, err
		}
		if channelInfo.PropertyId != hexJsonObj.Get("propertyid").Int() {
			err = errors.New("wrong propertyId at payerRsmcRdHex  at 41 protocol")
			log.Println(err)
			return nil, true, err
		}
		if commitmentTransaction.AmountToRSMC != hexJsonObj.Get("amount").Float() {
			err = errors.New("wrong amount at payerRsmcRdHex  at 41 protocol")
			log.Println(err)
			return nil, true, err
		}

		err = saveRdTx(tx, &channelInfo, signedRsmcHex, payerRsmcRdHex, commitmentTransaction, payerAddress, &user)
		if err != nil {
			return nil, false, err
		}
	}
	//endregion

	//region 5、对ht1a进行二次签名，并保存
	payerHt1aHex := jsonObj.C3aHtlcHtPartialSignedData.Hex
	payerHt1aInputsFromHtlc, err := getInputsForNextTxByParseTxHashVout(
		signedHtlcHex,
		commitmentTransaction.HTLCMultiAddress,
		commitmentTransaction.HTLCMultiAddressScriptPubKey,
		commitmentTransaction.HTLCRedeemScript)
	if err != nil {
		log.Println(err)
		return nil, true, err
	}
	multiAddress, _, _, err := createMultiSig(htlcRequestInfo.CurrHtlcTempAddressForHt1aPubKey, payeePubKey)
	if err != nil {
		log.Println(err)
		return nil, false, err
	}

	result, err = rpcClient.OmniDecodeTransactionWithPrevTxs(payerHt1aHex, payerHt1aInputsFromHtlc)
	if err != nil {
		log.Println(err)
		return nil, true, err
	}
	hexJsonObj = gjson.Parse(result)
	if commitmentTransaction.HTLCMultiAddress != hexJsonObj.Get("sendingaddress").String() {
		err = errors.New("wrong inputAddress at payerHt1aHex  at 41 protocol")
		log.Println(err)
		return nil, true, err
	}

	if multiAddress != hexJsonObj.Get("referenceaddress").String() {
		err = errors.New("wrong outputAddress at payerHt1aHex  at 41 protocol")
		log.Println(err)
		return nil, true, err
	}
	if channelInfo.PropertyId != hexJsonObj.Get("propertyid").Int() {
		err = errors.New("wrong propertyId at payerHt1aHex  at 41 protocol")
		log.Println(err)
		return nil, true, err
	}
	if commitmentTransaction.AmountToHtlc != hexJsonObj.Get("amount").Float() {
		err = errors.New("wrong amount at payerHt1aHex  at 41 protocol")
		log.Println(err)
		return nil, true, err
	}
	htlcTimeOut := commitmentTransaction.HtlcCltvExpiry
	ht1a, err := saveHT1aForAlice(tx, channelInfo, commitmentTransaction, payerHt1aHex, htlcRequestInfo, payeePubKey, htlcTimeOut, user)
	if err != nil {
		err = errors.New("fail to sign  payerHt1aHex  at 41 protocol")
		log.Println(err)
		return nil, true, err
	}

	//endregion

	//region 6、为bob存储lockByHForBobHex
	lockByHForBobHex := jsonObj.C3aHtlcHlockPartialSignedData.Hex
	if tool.CheckIsString(&lockByHForBobHex) == false {
		err = errors.New("lockByHForBobHex is empty at 41 protocol")
		log.Println(err)
		return nil, true, err
	}
	_, err = saveHtlcLockByHTxAtPayerSide(tx, channelInfo, commitmentTransaction, lockByHForBobHex, user)
	if err != nil {
		err = errors.New("fail to lockByHForBobHex at 41 protocol")
		log.Println(err)
		return nil, false, err
	}
	//endregion
	return &ht1a.RSMCTxHex, true, nil
}

// 收款方更新C3b的信息
func checkHexAndUpdateC3bOn42Protocal(tx storm.Node, jsonObj bean.NeedBobSignHtlcSubTxOfC3bP2p, channelInfo dao.ChannelInfo, latestCommitmentTx *dao.CommitmentTransaction, user bean.User) (data *dao.CommitmentTransaction, needNoticePayee bool, err error) {
	bobChannelAddress := channelInfo.AddressB
	aliceChannelAddress := channelInfo.AddressA
	if user.PeerId == channelInfo.PeerIdA {
		aliceChannelAddress = channelInfo.AddressB
		bobChannelAddress = channelInfo.AddressA
	}
	//region 1、检测 signedToOtherHex
	signedToCounterpartyTxHex := jsonObj.C3bCompleteSignedCounterpartyHex
	if tool.CheckIsString(&signedToCounterpartyTxHex) {
		_, err = rpcClient.TestMemPoolAccept(signedToCounterpartyTxHex)
		if err != nil {
			err = errors.New("wrong signedToCounterpartyTxHex at 42 protocol")
			log.Println(err)
			return nil, true, err
		}
		result, err := rpcClient.OmniDecodeTransaction(signedToCounterpartyTxHex)
		if err != nil {
			return nil, true, err
		}
		hexJsonObj := gjson.Parse(result)
		if channelInfo.ChannelAddress != hexJsonObj.Get("sendingaddress").String() {
			err = errors.New("wrong inputAddress at signedToOtherHex  at 42 protocol")
			log.Println(err)
			return nil, true, err
		}
		if channelInfo.PropertyId != hexJsonObj.Get("propertyid").Int() {
			err = errors.New("wrong propertyId at signedToOtherHex  at 42 protocol")
			log.Println(err)
			return nil, true, err
		}

		if aliceChannelAddress != hexJsonObj.Get("referenceaddress").String() {
			err = errors.New("wrong outputAddress at signedToOtherHex  at 42 protocol")
			log.Println(err)
			return nil, true, err
		}
		if latestCommitmentTx.AmountToCounterparty != hexJsonObj.Get("amount").Float() {
			err = errors.New("wrong amount at signedToOtherHex  at 42 protocol")
			log.Println(err)
			return nil, true, err
		}
		latestCommitmentTx.ToCounterpartyTxHex = signedToCounterpartyTxHex
		latestCommitmentTx.ToCounterpartyTxid = hexJsonObj.Get("txid").String()
	}
	//endregion

	//region 2、检测 signedRsmcHex
	signedRsmcHex := ""
	if len(latestCommitmentTx.RSMCTxHex) > 0 {
		signedRsmcHex = jsonObj.C3bCompleteSignedRsmcHex
		if pass, _ := rpcClient.CheckMultiSign(true, signedRsmcHex, 2); pass == false {
			err = errors.New("signedRsmcHex is empty at 42 protocol")
			log.Println(err)
			return nil, true, err
		}
		result, err := rpcClient.OmniDecodeTransaction(signedRsmcHex)
		if err != nil {
			return nil, true, err
		}
		hexJsonObj := gjson.Parse(result)

		if channelInfo.ChannelAddress != hexJsonObj.Get("sendingaddress").String() {
			err = errors.New("wrong inputAddress at signedRsmcHex  at 42 protocol")
			log.Println(err)
			return nil, true, err
		}
		if channelInfo.PropertyId != hexJsonObj.Get("propertyid").Int() {
			err = errors.New("wrong propertyId at signedRsmcHex  at 42 protocol")
			log.Println(err)
			return nil, true, err
		}

		if latestCommitmentTx.RSMCMultiAddress != hexJsonObj.Get("referenceaddress").String() {
			err = errors.New("wrong outputAddress at signedRsmcHex  at 42 protocol")
			log.Println(err)
			return nil, true, err
		}
		if latestCommitmentTx.AmountToRSMC != hexJsonObj.Get("amount").Float() {
			err = errors.New("wrong amount at signedRsmcHex  at 42 protocol")
			log.Println(err)
			return nil, true, err
		}
		latestCommitmentTx.RSMCTxHex = signedRsmcHex
		latestCommitmentTx.RSMCTxid = hexJsonObj.Get("txid").String()
	}
	//endregion

	//region 3、检测 signedHtlcHex
	signedHtlcHex := jsonObj.C3bCompleteSignedHtlcHex
	if pass, _ := rpcClient.CheckMultiSign(true, signedHtlcHex, 2); pass == false {
		err = errors.New("signedHtlcHex is empty at 42 protocol")
		log.Println(err)
		return nil, true, err
	}
	result, err := rpcClient.OmniDecodeTransaction(signedHtlcHex)
	if err != nil {
		return nil, true, err
	}
	hexJsonObj := gjson.Parse(result)
	if channelInfo.ChannelAddress != hexJsonObj.Get("sendingaddress").String() {
		err = errors.New("wrong inputAddress at signedHtlcHex  at 42 protocol")
		log.Println(err)
		return nil, true, err
	}
	if channelInfo.PropertyId != hexJsonObj.Get("propertyid").Int() {
		err = errors.New("wrong propertyId at signedHtlcHex  at 42 protocol")
		log.Println(err)
		return nil, true, err
	}

	if latestCommitmentTx.HTLCMultiAddress != hexJsonObj.Get("referenceaddress").String() {
		err = errors.New("wrong outputAddress at signedHtlcHex  at 42 protocol")
		log.Println(err)
		return nil, true, err
	}
	if latestCommitmentTx.AmountToHtlc != hexJsonObj.Get("amount").Float() {
		err = errors.New("wrong amount at signedHtlcHex  at 42 protocol")
		log.Println(err)
		return nil, true, err
	}
	latestCommitmentTx.HtlcTxHex = signedHtlcHex
	latestCommitmentTx.HTLCTxid = hexJsonObj.Get("txid").String()
	//endregion

	//region 4、rsmc Rd
	if len(latestCommitmentTx.RSMCTxHex) > 0 {
		payeeRsmcRdHex := jsonObj.C3bRsmcRdPartialData.Hex
		if pass, _ := rpcClient.CheckMultiSign(false, payeeRsmcRdHex, 2); pass == false {
			err = errors.New("signedRsmcHex is empty at 42 protocol")
			log.Println(err)
			return nil, true, err
		}
		payerRDInputsFromRsmc, err := getInputsForNextTxByParseTxHashVout(
			signedRsmcHex,
			latestCommitmentTx.RSMCMultiAddress,
			latestCommitmentTx.RSMCMultiAddressScriptPubKey,
			latestCommitmentTx.RSMCRedeemScript)
		if err != nil {
			log.Println(err)
			return nil, true, err
		}
		result, err = rpcClient.OmniDecodeTransactionWithPrevTxs(payeeRsmcRdHex, payerRDInputsFromRsmc)
		if err != nil {
			log.Println(err)
			return nil, true, err
		}
		hexJsonObj = gjson.Parse(result)
		if latestCommitmentTx.RSMCMultiAddress != hexJsonObj.Get("sendingaddress").String() {
			err = errors.New("wrong inputAddress at payeeRsmcRdHex  at 41 protocol")
			log.Println(err)
			return nil, true, err
		}

		if bobChannelAddress != hexJsonObj.Get("referenceaddress").String() {
			err = errors.New("wrong outputAddress at payeeRsmcRdHex  at 41 protocol")
			log.Println(err)
			return nil, true, err
		}
		if channelInfo.PropertyId != hexJsonObj.Get("propertyid").Int() {
			err = errors.New("wrong propertyId at payeeRsmcRdHex  at 41 protocol")
			log.Println(err)
			return nil, true, err
		}
		if latestCommitmentTx.AmountToRSMC != hexJsonObj.Get("amount").Float() {
			err = errors.New("wrong amount at payeeRsmcRdHex  at 41 protocol")
			log.Println(err)
			return nil, true, err
		}

		err = saveRdTx(tx, &channelInfo, signedRsmcHex, payeeRsmcRdHex, latestCommitmentTx, bobChannelAddress, &user)
		if err != nil {
			return nil, false, err
		}
	}
	//endregion

	// region  5 保存Hlock h+bobChannelPubkey锁定给bob的付款金额 有了H对应的R，就能解锁
	lockByHForBobHex := jsonObj.C3bHtlcHlockPartialData.Hex
	if tool.CheckIsString(&lockByHForBobHex) == false {
		err = errors.New("payeeHlockHex is empty at 41 protocol")
		log.Println(err)
		return nil, true, err
	}
	_, err = saveHtlcLockByHForBobAtPayeeSide(tx, channelInfo, latestCommitmentTx, lockByHForBobHex, user)
	if err != nil {
		err = errors.New("fail to signHtlcLockByHTxAtPayerSide at 41 protocol")
		log.Println(err)
		return nil, false, err
	}
	//endregion

	// region  6、签名HTD1b 超时退给alice的钱
	payeeHTD1bHex := jsonObj.C3bHtlcHtdPartialData.Hex
	if tool.CheckIsString(&payeeHTD1bHex) == false {
		err = errors.New("payeeHTD1bHex is empty at 41 protocol")
		log.Println(err)
		return nil, true, err
	}
	payeeHTD1bInputsFromHtlc, err := getInputsForNextTxByParseTxHashVout(
		signedHtlcHex,
		latestCommitmentTx.HTLCMultiAddress,
		latestCommitmentTx.HTLCMultiAddressScriptPubKey,
		latestCommitmentTx.HTLCRedeemScript)
	if err != nil {
		log.Println(err)
		return nil, true, err
	}
	result, err = rpcClient.OmniDecodeTransactionWithPrevTxs(payeeHTD1bHex, payeeHTD1bInputsFromHtlc)
	if err != nil {
		log.Println(err)
		return nil, true, err
	}
	hexJsonObj = gjson.Parse(result)
	if latestCommitmentTx.HTLCMultiAddress != hexJsonObj.Get("sendingaddress").String() {
		err = errors.New("wrong inputAddress at payeeHTD1bHex  at 41 protocol")
		log.Println(err)
		return nil, true, err
	}
	if aliceChannelAddress != hexJsonObj.Get("referenceaddress").String() {
		err = errors.New("wrong outputAddress at payeeHTD1bHex  at 41 protocol")
		log.Println(err)
		return nil, true, err
	}
	if channelInfo.PropertyId != hexJsonObj.Get("propertyid").Int() {
		err = errors.New("wrong propertyId at payeeHTD1bHex  at 41 protocol")
		log.Println(err)
		return nil, true, err
	}
	if latestCommitmentTx.AmountToHtlc != hexJsonObj.Get("amount").Float() {
		err = errors.New("wrong amount at payeeHTD1bHex  at 41 protocol")
		log.Println(err)
		return nil, true, err
	}

	err = saveHTD1bTx(tx, signedHtlcHex, payeeHTD1bHex, *latestCommitmentTx, aliceChannelAddress, &user)
	if err != nil {
		return nil, false, err
	}
	//endregion

	latestCommitmentTx.RsmcRawTxData = bean.NeedClientSignTxData{}
	latestCommitmentTx.HtlcRawTxData = bean.NeedClientSignTxData{}
	latestCommitmentTx.ToCounterpartyRawTxData = bean.NeedClientSignTxData{}

	latestCommitmentTx.CurrState = dao.TxInfoState_Htlc_GetH
	bytes, err := json.Marshal(latestCommitmentTx)
	msgHash := tool.SignMsgWithSha256(bytes)
	latestCommitmentTx.CurrHash = msgHash
	_ = tx.Update(latestCommitmentTx)

	return latestCommitmentTx, false, nil
}
