
function alertTrade(stockId, amount, response){
	consol.log("your order of stock ID: amount: was :")
} 




responses = {}
requestID = 0

function OnMessage(msg){
	if msg.action == "response"{
		responses[msg.request](msg)
		delete responses[msg.request]
	}
}

function MakeRequst(msg, callback){
	r = requestID
	responses[requestID++] = callback
	Websocket.Send(msg)
	return r
}

function MakeTrade(tradeMessage){
	MakeRequst(tradeMessage, function(response){
		alertTrade(tradeMessage.stockId, tradeMessage.amount, response)
	})
}

function SetGraphData(GraphName, DataSeries, QueryResponse){
	...
}

function GetPortfolioHistoricalData(portfolioUUID){
	requestIds = []
	GraphData = {}
	GraphData["GraphName"]	// request Networh
	requestIds.append(MakeRequst({Networth_Request}, function(response){
		GraphData["GraphName"]["networth"] = processQueryResponse(response)
	}))
	requestIds.append(MakeRequst({Wallet_Request}, function(response){
		GraphData["GraphName"]["wallet"] = processQueryResponse(response)

	}))
	for stock in stocksForPortfolio(portfolioUUID){
		requestIds.append((MakeRequst({stock_price_value}, function(response){
			GraphData["GraphName"]["stock_"+stock.uuid] = processQueryResponse(response)
		})
	}	
	var drawGraphOnceDone = null

	drawGraphOnceDone = function(){
		stillWaiting = false
		for(r in requestIds) {
			if(r in responses) {
				stillWaiting = true
			}
		}
		if !stillWaiting {
			DrawGraph(GraphData)
		}
		else{
			Set_Timeout(drawGraphOnceDone, 100)
		}
	}

	Set_Timout(drawGraphOnceDone, 100)

}

