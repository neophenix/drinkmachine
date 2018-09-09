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

        this.pour_drink = function(id, override) {
            // Start out hiding all the various buttons, info divs, etc
            $("#pour-drink").hide();
            $("#pour-drink-missing").hide();
            $("#drink-info").hide();
            $("#dispensing").hide();
            $("#missing").hide();
            $("#finish").hide();

            // We will only show the finishing items if we have any
            var show_finish = false;

            // set to -1 so we know when we get our first time_remaining message
            var time_remaining = -1;
            var ws = new WebSocket("ws://"+window.location.host+"/ws");
            ws.onopen = function(e) {
                ws.send(JSON.stringify({"action": "make_drink", "id": parseInt(id, 10), "options": {"override": override}}));
                // reset drink info, ingredients, approx time, etc
                $("#dispensing").html("Dispensing:<br/>");
                $("#finish").html("Finish with:<br/>");
            }
            ws.onmessage = function(msg) {
                msg = JSON.parse(msg.data);
                console.log(msg);
                switch (msg.type) {
                    case 'error':
                        DrinkMachine.show_alert("danger", msg.message);
                        $("#drink-info").hide();
                        $("#dispensing").hide();
                        $("#finish").hide();
                        ws.close();
                        break;
                    case 'starting':
                        $("#drink-info").show();
                        $("#dispensing").show();
                        break;
                    case 'pouring':
                        $("#dispensing").append($("<div>").attr("id", msg.id).text(msg.message));
                        $("#stop").show();
                        break;
                    case 'pour_complete':
                        $("#dispensing").find("[id='"+msg.id+"']").remove();
                        break;
                    case 'missing':
                        $("#pour-drink-missing").show();
                        $("#missing").html("Missing Ingredients:<br/>" + msg.message).show();
                        break;
                    case 'finish':
                        show_finish = true;
                        $("#finish").append($("<div>").text(msg.message));
                        break;
                    case 'notes':
                        show_finish = true;
                        $("#finish").append($("<div>").text(msg.message));
                        break;
                    case 'time_remaining':
                        var backend_time_remaining = parseFloat(msg.message);
                        if (backend_time_remaining == 0) {
                            $("#progress").hide();
                            $("#dispensing").hide();
                            if (show_finish) {
                                $("#finish").show();
                            }
                            ws.close();
                        }
                        else if (time_remaining == -1) {
                            time_remaining = backend_time_remaining;
                            $("#progress").show();
                            $("#drink-timer").attr("aria-valuemax", time_remaining)
                            setTimeout(function() { DrinkMachine.countdown(time_remaining, time_remaining); },0);
                        }
                        break;
                }
            }
            ws.onclose = function() {
                DrinkMachine.pages.pour_drink.finish_drink();
            }
        };

        this.run_pump = function(id, duration) {
            var ws = new WebSocket("ws://"+window.location.host+"/ws");
            ws.onopen = function(e) {
                ws.send(JSON.stringify({"action": "run_pump", "id": id, "options": { "duration": duration } }));
            }
            ws.onmessage = function(msg) {
                msg = JSON.parse(msg.data);
                switch (msg.type) {
                    case 'done':
                        DrinkMachine.show_alert("success", "Done");
                        ws.close();
                        break;
                    case 'error':
                        DrinkMachine.show_alert("danger", msg.message);
                        ws.close();
                        break;
                }
            }
        };

        this.stop_pumps = function() {
            var ws = new WebSocket("ws://"+window.location.host+"/ws");
            ws.onopen = function(e) {
                ws.send(JSON.stringify({"action": "stop", "id": 0}));
            }
        };
    };

    this.countdown = function(total, remaining) {
        var pct = (remaining / total) * 100;
        $("#drink-timer").attr("aria-valuenow",remaining).width(pct+"%");
        if (remaining > 0) {
            setTimeout(function() { DrinkMachine.countdown(total, remaining-1); }, 1000);
        }
    };

    var alert_last_type = "success"; // -_-
    this.show_alert = function(type, message) {
        $("#alert-msg").html(message);
        $("#alert").removeClass("alert-"+alert_last_type).addClass("alert-"+type).show();
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
