var vm_research;
var graphType, graphUsers, graphStocks;

function load_research_tab() {
    vm_research = new Vue({
        el: "#research--view",
        data: {
            receipt: {
                ticker: "",
                time: new Date(),
            }
        },
        methods: {
            drawGraph: function() {
                // Get user set variables for graphing
                var fields = [];
                var uuids = [];
                var type = getSelectized('#research-graph-type-select')[0];
                if (type === "stock") {
                    let selected = getSelectized('#research-graph-stock-select');
                    selected.forEach(function(d) {
                        uuids.push(d);
                        fields.push('current_price');
                    })
                } else if (type === "portfolio") {
                    let selected = getSelectized('#research-graph-user-select');
                    selected.forEach(function(d) {
                        uuids.push(d);
                        uuids.push(d);
                        fields.push('net_worth');
                        fields.push('wallet');
                    })
                }

                // Create the graph
                queryDrawGraph("#research-graph-svg-main", uuids, fields, false, false);

            },
            // updateSelections: function() {
            //     console.log("herer")
            //     var type = getSelectized('#research-graph-type-select')[0];
            //     if (type === 'stock') {
            //         $('#research-graph-user-select').hide();
            //         $('#research-graph-stock-select').show();
            //     } else if (type === 'portfolio') {
            //         $('#research-graph-stock-select').hide();
            //         $('#research-graph-user-select').show();
            //     }
            // }
            openTradeHistory: function(uuid) {
                var trade = this.tradeHistory.filter(d => d.uuid === uuid)[0];

                
                console.log(trade);
                
                this.receipt.ticker = vm_stocks.stocks[trade.stock_uuid].ticker_id;
                this.receipt.time = trade.time;

            }
        },
        computed: {
            tradeHistory: function() {
                var entries = Object.values(vm_recordEntry.entries)
                    .map(function(d) {
                        let record = vm_recordBook.records[d.book_uuid];
                        d.time = Date(d.time);
                        d.portfolio_uuid = record.portfolio_uuid;
                        d.ledger_uuid = record.ledger_uuid;
                        d.stock_uuid = record.stock_uuid;
                        return d; 
                    });
                    // entries = entries.filter(d => d.portfolio_uuid === )
                    return entries;
            }
        }
    });

    // Create selectize areas
    graphType = $('#research-graph-type-select').selectize({maxItems: 1});
    graphStocks = $('#research-graph-stock-select').selectize({maxItems: 5});
    graphUsers = $('#research-graph-user-select').selectize({maxItems: 5});
    // Start with users selection hidden
    $('#research-graph-user-select').hide();

};

function updateResearchStocks() {
    // Get the html element to update
    var $select = $(document.getElementById('research-graph-stock-select'));
    var selectize = $select[0].selectize;

    Object.values(vm_stocks.stocks)
        .map(function(d) {
            selectize.addOption({
                value: d.uuid,
                text: d.ticker_id
            });
        });
}

function updateResearchUsers() {
    // Get the html element to update
    var $select = $(document.getElementById('research-graph-user-select'));
    var selectize = $select[0].selectize;

    Object.values(vm_users.users)
        .forEach(function(d) {
            selectize.addOption({
                value: d.portfolio_uuid,
                text: d.display_name
            });
        });
}