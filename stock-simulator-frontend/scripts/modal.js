var buySellModal;
var transferModal;
var genericTextFieldModal;

function toggleModal() {
    $("#modal--container").toggleClass("open");
    // Set styling
    if (buySellModal.isBuying) {
        $("#calc-btn-buy").addClass("fill");
        $("#calc-btn-sell").removeClass("fill");
    }
    if (!buySellModal.isBuying) {
        $("#calc-btn-sell").addClass("fill");
        $("#calc-btn-buy").removeClass("fill");
    }
}

function toggleTransferModal() {
    console.log("Show generic modal");
    $("#transfer-Modal--container").toggleClass("open");
}

function toggleGenericTextFieldModal() {
    console.log("Show generic text field modal");
    $("#generic-text-field-modal--container").toggleClass("open");
}

function sendTrade(stock_id, amt) {
    if (amt != 0) {
        // Creating message for the trade request
        var msg = {
            stock_id: stock_id,
            amount: amt
        };
    
        var callback = function(msg) {
            var success = msg.msg.success;
    
            // If trade was a success
            if (success) {
                // Vars needed to form note
                var amount = Number(msg.msg.order.amount);
                console.log(msg.msg.order.stock_id)
                var stock_item = vm_stocks.stocks[msg.msg.order.stock_id];
                
                if (amount < 0) {
                    amount *= -1;
                    message = "Successful sale of " + amount + " " + stock_item.ticker_id + " stocks.";
                } else {
                    message = "Successful purchase of " + amount + " " + stock_item.ticker_id + " stocks."; 
                }
                notifyTopBar(message, GREEN, success);
    
            } else {
                message = msg.msg.err;
                notifyTopBar(message, RED, success);
            }
        };
    
        // Sending through websocket
        console.log("SEND TRADE");
        console.log(JSON.stringify(msg));
    
        // Send through WebSocket
        doSend("trade", msg, callback);
    
        // Reset buy sell amount
        buySellModal.buySellAmount = 0;

    }
}


function load_modal_vues() {

    // Vue object for the buy and sell modal
    buySellModal = new Vue({
        el: "#modal--container",
        data: {
            showModal: false,
            buySellAmount: 0,
            isBuying: true,
            stock_uuid: "Mockstarket",
            prospectiveCash: 0,
            prospectiveBonus: 0,
            prospectiveFees: 0,
            prospectiveShareCount: 0,
            prospectiveTax: 0,
            prospectiveResult: 0,
        },
        methods: {
            toPrice: formatPrice,
            addAmount: function(amt) {
                buySellModal.buySellAmount += amt;
                $('#buy-sell-amount-input').val(buySellModal.buySellAmount);
                $('#buy-sell-amount-input').focus();
            },
            setAmount: function(evt) {
                var user_input = $('#buy-sell-amount-input').val();
                if (!isNaN(user_input)) {
                    console.log(user_input)
                    buySellModal.buySellAmount = Number(user_input);
                }
            },
            clearAmount: function() {
                buySellModal.buySellAmount = 0;
                $('#buy-sell-amount-input').val(buySellModal.buySellAmount);
                $('#buy-sell-amount-input').focus();
            },
            setMax: function() {

                buySellModal.buySellAmount = 1000;
                

                // var user_max;
                
                // if (buySellModal.isBuying) {
                //     var user_level = vm_portfolios.currentUser.level;
                //     user_max = vm_store.levels.filter(d => d.level == user_level)[0].max_shares;

                //     var stock = vm_dash_tab.currUserStocks.filter(
                //         d => d.stock_id === buySellModal.stock_uuid
                //         )[0];
                //     var user_holdings = stock.amount;
                    
                //     var purchase_max = user_max - user_holdings;
                    
                //     buySellModal.buySellAmount = purchase_max;
                //     $('#buy-sell-amount-input').val(purchase_max);
                //     $('#buy-sell-amount-input').focus();

                // } else {
                //     //determine current users holdings
                //     var stock = vm_dash_tab.currUserStocks.filter(
                //         d => d.stock_id === buySellModal.stock_uuid
                //     )[0];
                //         console.log(stock)
                //     if (stock !== undefined) {
                //         user_max = stock.amount;
                //         if (buySellModal.buySellAmount > user_max) {
                //             buySellModal.buySellAmount = user_max;
                //             console.log(user_max)
                //             $('#buy-sell-amount-input').val(user_max);
                //             // $('#buy-sell-amount-input').focus();
                //         } else {
                //             buySellModal.buySellAmount = user_max;
                //             $('#buy-sell-amount-input').val(user_max);
                //         }
                //     } else {
                //         buySellModal.buySellAmount = 0;
                //         $('#buy-sell-amount-input').val(buySellModal.buySellAmount);
                //         $('#buy-sell-amount-input').focus();
                //     }
                // }
            },
            controlMax: function() {
                
                var user_level = vm_portfolios.currentUser.level;

                var user_max_shares = vm_store.levels.filter(d => d.level == user_level)[0].max_shares;
                
                var user_holdings = 0;
                try {
                    user_holdings = vm_dash_tab.currUserStocks.filter(
                        d => d.stock_id === buySellModal.stock_uuid
                    )[0].amount;
                } catch (err) {

                }
                
                // When user is buying
                if (buySellModal.isBuying) {
                    
                    var purchase_max = user_max_shares - user_holdings;

                    // When modal amount is greater than possible holdings
                    if (buySellModal.buySellAmount > purchase_max) {
                        buySellModal.buySellAmount = purchase_max;
                        $('#buy-sell-amount-input').val(purchase_max);
                        $('#buy-sell-amount-input').focus();
                    }
                    
                    // TODO: Control max based on prospectiveTrade return
                    
                // When user is selling
                } else {
                    // When selling more than is owned
                    if (buySellModal.buySellAmount > user_holdings) {
                        buySellModal.buySellAmount = user_holdings;
                        $('#buy-sell-amount-input').val(user_holdings);
                        $('#buy-sell-amount-input').focus();
                    }
                }
            },
            setIsBuying: function(bool) {
                console.log("IS_BUYING: "+bool);
                // Change buying or selling
                buySellModal.isBuying = bool;

                // Set styling
                if (buySellModal.isBuying) {
                    $("#calc-btn-buy").addClass("fill");
                    $("#calc-btn-sell").removeClass("fill");
                } 
                if (!buySellModal.isBuying) {
                    $("#calc-btn-sell").addClass("fill");
                    $("#calc-btn-buy").removeClass("fill");
                }
            },
            submitTrade: function() {
                // Change amount depending on buy/sell
                if (!buySellModal.isBuying) {
                    buySellModal.buySellAmount *= -1;
                }
                $('#buy-sell-amount-input').blur();
                $('#buy-sell-amount-input').val("");
                sendTrade(buySellModal.stock_uuid, buySellModal.buySellAmount);
                toggleModal();
            },
            closeModal: function() {
                // $('#buy-sell-amount-input').val(0);
                $('#buy-sell-amount-input').blur();
                $('#buy-sell-amount-input').val("");
                toggleModal();
                buySellModal.buySellAmount = 0;
                buySellModal.showModal = false;
                buySellModal.isBuying = true;
            }
        },
        computed: {
            stock: function() {
                var clickedStock = Object.values(vm_stocks.stocks).filter(
                d => d.uuid === buySellModal.stock_uuid
                )[0];
                return clickedStock;
            },
            user: function() {
                var currUserUUID = sessionStorage.getItem("uuid");
                if (vm_users.users[currUserUUID] !== undefined) {
                var currUserFolioUUID = vm_users.users[currUserUUID].portfolio_uuid;
                if (vm_portfolios.portfolios[currUserFolioUUID] !== undefined) {
                    var folio = vm_portfolios.portfolios[currUserFolioUUID];
                    folio.investments = folio.net_worth - folio.wallet;
                    return folio;
                }
                }
                return {};
            }
        },
        watch: {
            // Resetting amount if more than can be traded is selected
            buySellAmount: function() {
                // if (buySellModal.isBuying) {
                //     if (buySellModal.buySellAmount > buySellModal.stock.open_shares) {
                //         buySellModal.buySellAmount = buySellModal.stock.open_shares;
                //     }
                //     // determine users cash and limit on purchase cost
                //     let cash = buySellModal.user.wallet;
                //     let purchase_val =
                //         buySellModal.stock.current_price * buySellModal.buySellAmount;
                //     if (purchase_val > cash) {
                //         buySellModal.buySellAmount = Math.floor(
                //         cash / buySellModal.stock.current_price
                //         );
                //     }
                // } else {
                //     //determine current users holdings
                //     let stock = vm_dash_tab.currUserStocks.filter(
                //         d => d.stock_id == buySellModal.stock_uuid
                //     )[0];
                //     if (stock !== undefined) {
                //         if (buySellModal.buySellAmount > stock.amount) {
                //             buySellModal.buySellAmount = stock.amount;
                //         }
                //     }
                // }

                // find max and set it there
                this.controlMax();

                // do a prospectiveTrade
                if (buySellModal.isBuying) {
                    var amount = buySellModal.buySellAmount;
                } else {
                    var amount = buySellModal.buySellAmount * (-1);
                }
                var callback = function(msg) {
                    if (msg.msg.success) {
                        updateModalFromProspect(msg);
                    }
                };
                prospectiveTrade(buySellModal.stock_uuid, amount, callback);
            },
            isBuying: function() {
                if (this.isBuying) {
                    console.log("is buying");
                    $("#calc-btn-buy").addClass("fill");
                    $("#calc-btn-sell").removeClass("fill");
                } else {
                    console.log("is selling");
                    $("#calc-btn-sell").addClass("fill");
                    $("#calc-btn-buy").removeClass("fill");
                }
                this.buySellAmount = 0;
                $('#buy-sell-amount-input').val(buySellModal.buySellAmount);
                // do a prospectiveTrade
                if (buySellModal.isBuying) {
                    var amount = buySellModal.buySellAmount;
                } else {
                    var amount = buySellModal.buySellAmount * (-1);
                }
                var callback = function(msg) {
                    if (msg.msg.success) {
                        updateModalFromProspect(msg);
                    }
                };
                prospectiveTrade(buySellModal.stock_uuid, amount, callback);
            }
        }
    });

    // Vue object for the buy and sell modal
    transferModal = new Vue({
        el: "#transfer-Modal--container",
        data: {
            showModal: false,
            recipient_uuid: '',
            recipient_name: '',
        },
        methods: {
            submitTransfer: function() {
                // Get current amount
                let amt = Number($('#cash-transfer-amount').val());
                amt *= 100;

                // Creating message for the transfer
                var msg = {
                    amount: amt,
                    recipient: transferModal.recipient_uuid
                };

                // Send through WebSocket
                doSend("transfer", msg); 
                
                // Close the modal
                $('#cash-transfer-amount').val('');
                $('#cash-transfer-amount').blur();
                toggleTransferModal();
            },
            closeModal: function() {
                transferModal.showModal = false;
                toggleTransferModal();
            }
        },
        watch: {
            recipient_name: function() {
                console.log("WEATCHERS GOING")
                $('#cash-transfer-target').val(transferModal.recipient_name);
                console.log(transferModal.recipient_name)
                $('#cash-transfer-target').val();
            }
        }
    });

    // Vue object for the buy and sell modal
    genericTextFieldModal = new Vue({
        el: "#generic-text-field-modal--container",
        data: {
        showModal: false,
        // investor_uuid: '',
        investor_name: "DieselBeaver"
        },
        methods: {
        toPrice: formatPrice,

        closeModal: function() {
            toggleGenericTextFieldModal();
            // transferModal.investor_uuid = '';
            // transferModal.investor_name = '';
            genericTextFieldModal.showModal = false;
        }
        },
        computed: {},
        watch: {}
    });

    $(".mini-calc-btn").click(function(event) {
        console.log("clicked mini btn");    
    });
}