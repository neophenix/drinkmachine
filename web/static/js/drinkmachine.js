var DrinkMachine = new function () {
    this.api = new function() {
        this.get = function(type, id = 0, cb) {
            $.get("/api/" + type + "/" + id, function(data) {
                cb(data);
            }, "json");
        };

        this.create = function(type, form, cb) {
            form = JSON.stringify(form)
            $.post("/api/" + type, form, function(data) {
                cb(data);
            }, "json");
        };

        this.update = function(type, id, form, cb) {
            form = JSON.stringify(form)
            $.ajax({
                url:        "/api/" + type + "/" + id,
                method:     "PUT",
                data:       form,
                dataType:   "json",
                success:    function(data) { cb(data); }
            });
        };

        this.remove = function(type, id, cb) {
            $.ajax({
                url:        "/api/" + type + "/" + id,
                method:     "DELETE",
                success:    function(data) { cb(data); }
            });
        };

        this.pour_drink = function(id) {
            $("#drink-info").hide();
            $("#dispensing").hide();
            $("#finish").hide();
            var time_remaining = -1;
            var ws = new WebSocket("ws://"+window.location.host+"/ws");
            ws.onopen = function(e) {
                ws.send(JSON.stringify({"action": "make_drink", "id": parseInt(id, 10)}));
                $("#drink-info").show();
                // drink info, ingredients, approx time, etc
                $("#dispensing").html("Dispensing:<br/>").show();
                $("#finish").html("Finish with:<br/>");
            }
            ws.onmessage = function(msg) {
                msg = JSON.parse(msg.data);
                console.log(msg);
                switch (msg.type) {
                    case 'pouring':
                        $("#dispensing").append($("<div>").attr("id", msg.id).text(msg.message));
                        break;
                    case 'pour_complete':
                        $("#dispensing").find("[id='"+msg.id+"']").remove();
                        break;
                    case 'finish':
                        $("#finish").append($("<div>").text(msg.message));
                        break;
                    case 'notes':
                        $("#finish").append($("<div>").text(msg.message));
                        break;
                    case 'time_remaining':
                        if (time_remaining == -1) {
                            time_remaining = parseFloat(msg.message)
                            $("#progress").show();
                            $("#drink-timer").attr("aria-valuemax", time_remaining)
                            setTimeout(function() { DrinkMachine.countdown(time_remaining, time_remaining); },0);
                        }
                        else {
                            time_remaining = parseFloat(msg.message)
                            if (time_remaining == 0) {
                                ws.close();
                            }
                        }
                        break;
                }
            }
            ws.onclose = function() {
                DrinkMachine.pages.pour_drink.finish_drink();
            }
        };
    };

    this.countdown = function(total, remaining) {
        var pct = (remaining / total) * 100;
        $("#drink-timer").attr("aria-valuenow",remaining).width(pct+"%");
        if (remaining > 0) {
            setTimeout(function() { DrinkMachine.countdown(total, remaining-1); }, 1000);
        }
        else {
            $("#progress").hide();
            $("#dispensing").hide();
            $("#finish").show();
        }
    };

    var alert_last_type = "success"; // -_-
    this.show_alert = function(type, message) {
        $("#alert-msg").html(message);
        $("#alert").removeClass(alert_last_type).addClass("alert-"+type).show();
        alert_last_type = type;
        setTimeout(DrinkMachine.hide_alert, 3000);
    };

    this.hide_alert = function() {
        $("#alert").hide();
    };
};

// just some basic stuff shared on each page
$("#admin-menu-btn").click(function() {
    $("#admin-menu").toggle();
});

$("#home-btn").click(function() {
    window.location.href = "/"; 
});
