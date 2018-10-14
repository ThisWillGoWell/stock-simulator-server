var vm_research;
var graphType, graphUsers, graphStocks;

function load_research_tab() {
    vm_research = new Vue({
        el: "#research--view",
        methods: {
            drawGraph: function() {
                console.log(getSelectized('#research-graph-type-select'))
                DrawLingGraph("#research-graph-svg-main", getGraphingData);
            },
            // updateSelections: function() {
            //     // Show needed select boxes
            //     // $('#research-graph-user-select').hide();
            //     // Remove unneeded selection boxes
            // }
        }
    });

    graphType = $('#research-graph-type-select').selectize({maxItems: 1});
    // $('#research-graph-stock-select').selectize({maxItems: 5});
    // $('#research-graph-user-select').selectize({maxItems: 5});
};

function updateResearchStocks() {
    var stocks = Object.values(vm_stocks.stocks)
        .map(function(d) {
            return {
                uuid: d.uuid,
                ticker_id: d.ticker_id
            };
        });
    console.log(stocks)
    var options = d3.select('#research-graph-stock-select').selectAll('option')
        .data(stocks, d => d.uuid);

    options.exit().remove();

    options.enter().append('option')
        .attr('value', d => d.uuid)
        .text(d => d.ticker_id);
}

function updateResearchUsers() {
    var users = Object.values(vm_users.users)
        .map(function(d) {
            return {
                uuid: d.portfolio_uuid,
                name: d.display_name
            };
        });

    var options = d3.select('#research-graph-user-select').selectAll('option')
        .data(users, d => d.portfolio_uuid);

    options.exit().remove();

    options.enter().append('option')
        .attr('value', d => d.uuid)
        .text(d => d.name);
}