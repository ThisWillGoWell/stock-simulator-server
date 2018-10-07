var buySellModal;
var transferModal;
var genericTextFieldModal;

function toggleModal() {
    $("#modal--container").toggleClass("open");
}

function toggleTransferModal() {
    console.log("Show generic modal");
    $("#transfer-Modal--container").toggleClass("open");
}

function toggleGenericTextFieldModal() {
    console.log("Show generic text field modal");
    $("#generic-text-field-modal--container").toggleClass("open");
}


// $("#modal--container").click(function() {
//     console.log("modal quit");
//     $("#modal--container").removeClass("open");
// });

function sendTrade() {
    // Creating message for the trade request
    var msg = {
        stock_id: buySellModal.stock_uuid,
        amount: buySellModal.buySellAmount
    };

    // Sending through websocket
    console.log("SEND TRADE");

    REQUESTS[REQUEST_ID] = function(msg) {
        if (msg.msg.success) {
            notify("Trade successful!", msg.msg.success);
        } else {
            notify("Trade unsuccessful: " + msg.msg.err, msg.msg.success);
        }
    };

    // Send through WebSocket
    console.log(JSON.stringify(msg));
    doSend("trade", msg);

    REQUEST_ID++;

    // Reset buy sell amount
    buySellModal.buySellAmount = 0;
}


function load_modal_vues() {

    // Vue object for the buy and sell modal
    buySellModal = new Vue({
        el: "#modal--container",
        data: {
        showModal: false,
        buySellAmount: 0,
        isBuying: true,
        stock_uuid: "OSRS"
        },
        methods: {
        toPrice: formatPrice,
        addAmount: function(amt) {
            buySellModal.buySellAmount += amt;
        },
        clearAmount: function() {
            buySellModal.buySellAmount = 0;
        },
        determineMax: function() {
            if (buySellModal.isBuying) {
            buySellModal.buySellAmount = buySellModal.stock.open_shares;
            } else {
            //determine current users holdings
            let stock = vm_dash_tab.currUserStocks.filter(
                d => d.stock_id === buySellModal.stock_uuid
            )[0];
            if (stock !== undefined) {
                buySellModal.buySellAmount = stock.amount;
            } else {
                buySellModal.buySellAmount = 0;
            }
            }
        },
        setIsBuying: function(bool) {
            // Change buying or selling
            buySellModal.isBuying = bool;

            // Set styling
            if (buySellModal.isBuying) {
            $("#calc-btn-buy").addClass("fill");
            $("#calc-btn-sell").removeClass("fill");
            } else {
            $("#calc-btn-sell").addClass("fill");
            $("#calc-btn-buy").removeClass("fill");
            }
        },
        submitTrade: function() {
            // Change amount depending on buy/sell
            if (!buySellModal.isBuying) {
            buySellModal.buySellAmount *= -1;
            }
            sendTrade();
            toggleModal();
        },
        closeModal: function() {
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
            if (buySellModal.isBuying) {
            if (buySellModal.buySellAmount > buySellModal.stock.open_shares) {
                buySellModal.buySellAmount = buySellModal.stock.open_shares;
            }
            // determine users cash and limit on purchase cost
            let cash = buySellModal.user.wallet;
            let purchase_val =
                buySellModal.stock.current_price * buySellModal.buySellAmount;
            if (purchase_val > cash) {
                buySellModal.buySellAmount = Math.floor(
                cash / buySellModal.stock.current_price
                );
            }
            } else {
            //determine current users holdings
            let stock = vm_dash_tab.currUserStocks.filter(
                d => d.stock_id == buySellModal.stock_uuid
            )[0];
            if (stock !== undefined) {
                if (buySellModal.buySellAmount > stock.amount) {
                buySellModal.buySellAmount = stock.amount;
                }
            }
            }
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
            console.log(amt);
            amt *= 100;

            // Creating message for the transfer
            var msg = {
            amount: amt,
            recipient: transferModal.recipient_uuid
            };

            REQUESTS[REQUEST_ID] = function(msg) {
            if (msg.msg.success) {
                notify("Transfer successful!", msg.msg.success);
            } else {
                notify("Transfer unsuccessful: " + msg.msg.err, msg.msg.success);
            }
            };

            // Send through WebSocket
            doSend("transfer", msg);

            REQUEST_ID++;

            // Close the modal
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
}