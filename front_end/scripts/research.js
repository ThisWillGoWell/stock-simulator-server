var vm_research;
var graphType, graphUsers, graphStocks;

function load_research_tab() {
    vm_research = new Vue({
        el: "#research--view",
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

                console.log(uuids)
                console.log(fields)
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
        },
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