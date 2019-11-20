// Send a prospective trade to see the result 
function prospectiveTrade(stock_id, amount, callback) {
    if (amount != 0) {
        var message = {
            stock_id: stock_id,
            amount: amount
        };
        
        var callback = callback;
    
        doSend("prospect", message, callback);
    }
};

function updateModalFromProspect(msg) {
    console.log(msg.msg.details.result);
    buySellModal.prospectiveCash = formatPrice(buySellModal.user.wallet + (msg.msg.details.result));
    buySellModal.prospectiveBonus = formatPrice(msg.msg.details.bonus)
    buySellModal.prospectiveFees = formatPrice(msg.msg.details.fees)
    buySellModal.prospectiveShareCount = msg.msg.details.share_count
    buySellModal.prospectiveTax = formatPrice(msg.msg.details.tax)
    buySellModal.prospectiveResult = formatPrice(msg.msg.details.result)
}



function prospectStockValues(stocks, realStockValue) {

    var newRealStockValue = 0;
    var stillWaiting = true;
	var responses = [];
    var requests = [];
    
    stocks.forEach(function(d) {
        if (d.amount == 0) {
            return 0;
        }
        requests.push(REQUEST_ID.toString());
        var callback = function(msg) {
            if (msg.msg.success) {
                newRealStockValue += Number(msg.msg.details.result);
            }
            
            responses.push(msg.request_id);
        };
        prospectiveTrade(d.stock_id, (-1)*d.amount, callback);
    })
    
  
	var whenDone = null;
	whenDone = function() {
	  if (requests.every(r => responses.indexOf(r) > -1)) {
		stillWaiting = false;
	  }
  
	  if (!stillWaiting) {
        // draw graph once all the data is back
        updateRealStocksValue(newRealStockValue);
	  } else {
		setTimeout(whenDone, 100);
	  }
	};
  
	setTimeout(whenDone, 100);
}

function updateRealStocksValue(new_value) {
    vm_dash_tab.currUserPortfolio.stocksValue = new_value;
    vm_dash_tab.currUserPortfolio.realNetworth = new_value + vm_dash_tab.currUserPortfolio.wallet;
}