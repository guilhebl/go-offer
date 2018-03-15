// loads table and attaches event listeners
$(document).ready(function () {
    // event handler for form
    $( "#addForm" ).submit(function( event ) {
        addOffer();
        event.preventDefault();
    });

    //event handler for reset btn
    $( "#resetBtn" ).click(function( event ) {
        resetList();
        event.preventDefault();
    });

    // event handler for tabs
    $( "#add-tab-btn" ).click(function( event ) {
        openTab(event, 'add-tab');
        event.preventDefault();
    });
    $( "#list-tab-btn" ).click(function( event ) {
        openTab(event, 'list-tab');
        event.preventDefault();
    });

    // open list tab
    openTab(null, 'list-tab');
});

function emptyList() {
    $('#tableBody').empty();
}

function loadList() {
    // load list
    $.ajax({
        type: "GET",
        url: "/offerlist",
        success: function(result) {
            emptyList();
            var list = $('#tableBody');
            for (var i in result.list)
            {
                var o = result.list[i];
                list.append(buildRow(o));
            }
        },
        error: function(e) {
            console.log("ERROR: e = " + e);
        }
    });
}

function resetList() {
    $.ajax({
        type: "GET",
        url: "/reset",
        success: function(result) {
            loadList();
        },
        error: function(e) {
            console.log("ERROR: e = " + e);
        }
    });
}

function buildRow(offer) {
    return "<tr>" +
            "<td>" + offer.id + "</td>" +
            "<td>" + offer.upc + "</td>" +
            "<td>" + offer.name + "</td>" +
            "<td>" + offer.partyName + "</td>" +
            "<td>" + offer.semanticName + "</td>" +
            "<td>" + offer.mainImageFileUrl + "</td>" +
            "<td>" + offer.partyImageFileUrl + "</td>" +
            "<td>" + offer.productCategory + "</td>" +
            "<td>" + offer.price + "</td>" +
            "<td>" + offer.rating + "</td>" +
            "<td>" + offer.numReviews + "</td>" +
           "</tr>";
}

// add Offer POST to REST API
function addOffer() {

    // build JSON
    var $form = $("#addForm");
    var o = {};
    o["upc"] = getFieldVal($form, 'upc');
    o["name"] = getFieldVal($form, 'name');
    o["partyName"] = getFieldVal($form, 'partyName');
    o["semanticName"] = getFieldVal($form, 'semanticName');
    o["mainImageFileUrl"] = getFieldVal($form, 'mainImageFileUrl');
    o["partyImageFileUrl"] = getFieldVal($form, 'partyImageFileUrl');
    o["productCategory"] = getFieldVal($form, 'productCategory');
    o["price"] = parseFloat(getFieldVal($form, 'price'));
    o["rating"] = parseFloat(getFieldVal($form, 'rating'));
    o["numReviews"] = parseInt(getFieldVal($form, 'numReviews'));
    var json = JSON.stringify(o);

    // post to REST Api
    $.ajax({
        type: "POST",
        url: "/offerlist",
        dataType : 'json', // data type
        data : json,
        success: function()
        {
            // open list tab
            openTab(null, 'list-tab');
        },
        error: function(e) {
            console.log("ERROR: e = " + e);
        }
    });
}

// opens a Tab
function openTab(evt, tabName) {
    openTabAction(tabName);

    var i, tabcontent, tablinks;
    tabcontent = document.getElementsByClassName("tabcontent");
    for (i = 0; i < tabcontent.length; i++) {
        tabcontent[i].style.display = "none";
    }
    tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(" active", "");
    }
    document.getElementById(tabName).style.display = "block";
    if (evt) evt.currentTarget.className += " active";
}

function openTabAction(tabName) {
    if (tabName == "list-tab") {
        // refresh list
        loadList();
    }
}

function getFieldVal($form, name) {
    return $form.find( "input[name='"+ name + "']" ).val()
}