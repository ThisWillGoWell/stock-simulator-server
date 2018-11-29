var vm_research;
var graphType, graphUsers, graphStocks, selectedStocks, selectedUsers;

function load_research_tab() {
    vm_research = new Vue({
        el: "#research--view",
        data: {
            receipt: {
                ticker: "",
                time: new Date(),
            },
            gType: "stocks",
        },
        methods: {
            drawGraph: function() {
                // Get user set variables for graphing
                var fields = [];
                var uuids = [];
                //var type = getSelectized('#research-graph-type-select')[0];
                var type = this.gType;
                console.log("ACTIVE CHART TYPE: "+type);
                if (type === "stocks") {
                    selectedStocks = getSelectized('#research-graph-stock-select');
                    selectedStocks.forEach(function(d) {
                        uuids.push(d);
                        fields.push('current_price');
                    })
                } else if (type === "investors") {
                    selectedUsers = getSelectized('#research-graph-user-select');
                    selectedUsers.forEach(function(d) {
                        uuids.push(d);
                        uuids.push(d);
                        fields.push('net_worth');
                        fields.push('wallet');
                    })
                }

                // Create the graph
                if (uuids.length > 0) {
                    queryDrawGraph("#research-graph-svg-main", uuids, fields, false, false);
                }

            },
            // updateSelections: function(newSelection) {
            //     console.log("herer");
            //     //var type = getSelectized('#research-graph-type-select')[0];
            //     if (newSelection === 'stock') {
            //         $('#query-term-users').addClass("shrunk");
            //         $('#query-term-users').removeClass("expanded");
            //         $('#query-term-stocks').addClass("expanded");
            //         $('#query-term-stocks').removeClass("shrunk");

            //     } else if (newSelection === 'portfolio') {
            //         $('#query-term-stocks').addClass("shrunk");
            //         $('#query-term-stocks').removeClass("expanded");
            //         $('#query-term-users').addClass("expanded");
            //         $('#query-term-users').removeClass("shrunk");
            //     }
            // },
            openTradeHistory: function(uuid) {
                var trade = this.tradeHistory.filter(d => d.uuid === uuid)[0];

                console.log(trade);
                
                this.receipt.ticker = vm_stocks.stocks[trade.stock_uuid].ticker_id;
                this.receipt.time = trade.time;

            },
            queryStocks: function() {
                if(this.gType !== "stocks") {
                    this.gType = "stocks";
                    $('#query-type-container .option.stocks').addClass("active");
                    $('#query-type-container .option.investors').removeClass("active");
                    $('.query-stocks').removeClass('hidden');
                    $('.query-investors').addClass('hidden');
                    $('#research-graph-svg-main').empty();
                    $('.query-items-label').text("STOCKS");
                    TweenMax.from($('.query-items-label'), 0.2, {ease: Back.easeOut.config(1.7), x:-75, opacity:0});
                    this.drawGraph();
                }
            },
            queryInvestors: function() {
                if(this.gType !== "investors") {
                    this.gType = "investors";
                    $('#query-type-container .option.investors').addClass("active");
                    $('#query-type-container .option.stocks').removeClass("active");
                    $('.query-stocks').addClass('hidden');
                    $('.query-investors').removeClass('hidden');
                    $('#research-graph-svg-main').empty();
                    $('.query-items-label').text("INVESTORS");
                    TweenMax.from($('.query-items-label'), 0.2, {ease: Back.easeOut.config(1.7), x:-75, opacity:0});
                    this.drawGraph();
                }
            },
            refreshGraph: function() {
                
                this.drawGraph();
                
            },
            
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

    function updateSelections(newSelection) {
        //console.log("herer");
        //var type = getSelectized('#research-graph-type-select')[0];
        if (newSelection === 'stock') {
            $('#query-term-users').addClass("shrunk");
            $('#query-term-users').removeClass("expanded");
            $('#query-term-stocks').addClass("expanded");
            $('#query-term-stocks').removeClass("shrunk");
            console.log("query stocks");
        } else if (newSelection === 'investors') {
            $('#query-term-stocks').addClass("shrunk");
            $('#query-term-stocks').removeClass("expanded");
            $('#query-term-users').addClass("expanded");
            $('#query-term-users').removeClass("shrunk");
            console.log("query investors");
        }
    }

    // Create selectize areas
    graphType = $('#research-graph-type-select').selectize({
        maxItems: 1,
        onChange: function(value) {
            updateSelections(value);
        }
    });
    graphStocks = $('#research-graph-stock-select').selectize(
        {
            maxItems: 5,
            onItemAdd(value, $item) {
                vm_research.drawGraph();
                //console.log($item);
            },
            onItemRemove(value) {
                if($('.has-items').length == 0) {
                    $('#research-graph-svg-main').empty();
                } else {
                    vm_research.drawGraph();
                }
            },
            
        }
    );
    graphUsers = $('#research-graph-user-select').selectize(
        {
            maxItems: 5,
            onItemAdd(value, $item) {
                vm_research.drawGraph();
                //console.log($item);
            },
            onItemRemove(value) {
                if($('.has-items').length == 0) {
                    $('#research-graph-svg-main').empty();
                } else {
                    vm_research.drawGraph();
                }
            },
            
        }
    );
    // Start with users selection hidden
    //$('#research-graph-user-select').hide();

    

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



// $("button.option.stocks").click(function(event) {
//     console.log("clicked stocks");    
// });

// $("button.option.investors").click(function(event) {
//     console.log("clicked investors");    
// });