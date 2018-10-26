function getUserRecordBook(userLedger) {
    var userRecordBook = Object.values(vm_recordBook.records)
        .filter(d => userLedger.indexOf(d.ledger_uuid) > -1)

    return userRecordBook;
}

function getUserRecordEntries(recordBook) {
    var entryUUIDs = [];
    recordBook.forEach(function (d) {
        d.buy_records.forEach(function(e) {
            console.log(e);
            entryUUIDs.push(e.RecordUuid);
        });
    })

    var entries = Object.values(vm_recordEntry.entries).filter(d => entryUUIDs.indexOf(d.uuid) > -1);

    return entries;
}

function getUserStockRecord(user_uuid, stock_id) {
    // Get user ledger entry that ties to this stock
    var userFolioUUID = vm_users.users[user_uuid].portfolio_uuid;
    var userLedger = Object.values(vm_ledger.ledger)
        .filter(d => d.portfolio_id === userFolioUUID)
        // .filter(d => d.stock_id === stock_id)
        .map(d => d.uuid);
    
    var userRecordBook = getUserRecordBook(userLedger);
    var userRecordEntries = getUserRecordEntries(userRecordBook);
    console.log(userRecordBook);
    console.log(userRecordEntries);
}