// Send a prospective trade to see the result 
function prospectiveTrade(stock_id, amount) {
    var message = {
        stock_id: stock_id,
        amount: amount
    };
    
    var callback = function(msg) {
        if (msg.msg.success) {
            updateModalFromProspect(msg);
        }
    };

    doSend("prospect", message, callback);
};

function updateModalFromProspect(msg) {
    console.log(msg.msg.details.result);
    buySellModal.prospectiveCash = formatPrice(buySellModal.user.wallet + (msg.msg.details.result));
}

