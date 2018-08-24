DrinkMachine.pages = new function() {
    this.pour_drink = new function() {
        this.init = function() {
            $("#drink").on("change", function() {
                if ($(this).val() == 0) {
                    $("#pour-drink").hide();
                }
                else {
                    $("#pour-drink").show();
                }
            });

            $("#pour-drink").on("click", function() {
                DrinkMachine.api.pour_drink($("#drink").val());
            });
        };

        this.finish_drink = function() {
            $("#drink").val(0);
            $("#pour-drink").hide();
        };
    };

    this.manage_drinks = new function() {
        this.init = function() {
            $("#drink").on("change", function() {
                if ($(this).val() == "_null") {
                    DrinkMachine.pages.manage_drinks.reset();
                }
                else if ($(this).val() == "_new") {
                    DrinkMachine.pages.manage_drinks.reset();
                    DrinkMachine.pages.manage_drinks.new_ingredient();
                    $("#name").show();
                    $("#notes").show();
                    $("#add-ingredient").show();
                    $("#add-drink").show();
                    $("#del-drink").hide();
                }
                else {
                    DrinkMachine.pages.manage_drinks.reset();
                    DrinkMachine.api.get("drink", $(this).val(), function(data) {
                        data.ingredients.forEach(function(ingredient) {
                            DrinkMachine.pages.manage_drinks.new_ingredient(ingredient);
                        });
                        $("#name").val(data.name).show();
                        $("#notes").val(data.notes).show();
                        $("#add-ingredient").show();
                        $("#add-drink").show();
                        $("#del-drink").show();
                    });
                }
            });

            $("#add-ingredient").on("click", function() {
                DrinkMachine.pages.manage_drinks.new_ingredient();
            });


            $("#add-drink").on("click", function() {
                if ($("#drink").val() == "_null") { return; }

                var form = {
                    name: $("#name").val(), 
                    notes: $("#notes").val(), 
                    ingredients: []
                }
                $(".ingredient-row").each(function(i,row) {
                    var ingredient = {
                        ingredient: $(row).find("[name^=ingredient-]").val(),
                        amount: parseFloat($(row).find("[name^=amount-]").val()),
                        units: $(row).find("[name^=units-]").val(),
                        dispense: $(row).find("[name^=dispense-]").prop("checked")
                    };

                    if ( ingredient.ingredient ) {
                        form.ingredients.push(ingredient);
                    }
                });

                if ($("#drink").val() !== "_new") {
                    // Updating existing drink
                    DrinkMachine.api.update("drink", $("#drink").val(), form, function(data) {
                        DrinkMachine.show_alert("success", "Your drink has been saved");
                        DrinkMachine.pages.manage_drinks.reset();
                    });
                }
                else {
                    // This is a new drink
                    DrinkMachine.api.create("drink", form, function(data) {
                        DrinkMachine.show_alert("success", "Your drink has been saved");
                        $("#drink").append($("<option>").val(data.id).text(data.name));
                        DrinkMachine.pages.manage_drinks.reset();
                    });
                }
            });

            $("#del-drink").on("click", function() {
                DrinkMachine.api.remove("drink", $("#drink").val(), function(data) {
                    DrinkMachine.show_alert("success", "Your drink has been removed");
                    $("#drink option[value='"+data.id+"']").remove();
                    DrinkMachine.pages.manage_drinks.reset();
                });
            });

            $("#ingredients-container").on("change", "[name^=ingredient-]", function() {
                if ($(this).val() != "" && $(this).data("count") == DrinkMachine.pages.manage_drinks.ingredient_count) {
                    DrinkMachine.pages.manage_drinks.new_ingredient();
                }
            });

            $("#ingredients-container").on("click", ".remove-ingredient", function() {
                $(this).parent().parent().parent().remove();
            });
        };

        this.ingredient_count = 0;
        this.new_ingredient = function(data = {units: "ml", dispense: 1, ingredient: ""}) {
            DrinkMachine.pages.manage_drinks.ingredient_count++;

            var $tmpl = $("#ingredient-row").clone();
            $tmpl.attr("id", "ingredient-row-" + DrinkMachine.pages.manage_drinks.ingredient_count);
            new_inputs($tmpl, "amount", data.amount);
            new_inputs($tmpl, "units", data.units);
            new_inputs($tmpl, "ingredient", data.ingredient);
            new_inputs($tmpl, "dispense", data.dispense);
            $tmpl.show();
            $("#ingredients-container").append($tmpl);

            function new_inputs ($tmpl, id, value) {
                $tmpl.find("#"+id)
                    .attr("name", id+"-" + DrinkMachine.pages.manage_drinks.ingredient_count)
                    .attr("id", id+"-" + DrinkMachine.pages.manage_drinks.ingredient_count)
                    .data("count", DrinkMachine.pages.manage_drinks.ingredient_count)
                    .val(value);

                if (id == "dispense" && value == true) {
                    $tmpl.find("#"+id+"-" + DrinkMachine.pages.manage_drinks.ingredient_count).attr("checked",true);
                }
            };
        };

        this.reset = function() {
            $("#ingredients-container").empty();
            DrinkMachine.pages.manage_drinks.ingredient_count = 0;
            $("#name").val("").hide();
            $("#notes").val("").hide();
            $("#add-ingredient").hide();
            $("#add-drink").hide();
            $("#del-drink").hide();
        };
    };

    this.manage_ingredients = new function() {
        this.init = function() {
            $("#ingredient-form").on("submit", function(e) {
                e.preventDefault();
            });

            $("#ingredient").on("change", function() {
                if ($(this).val() == "_new") {
                    $("#new-ingredient").show();
                    $("#add-ingredient").show();
                    $("#del-ingredient").hide();
                }
                else if ($(this).val() == "_null") {
                    $("#new-ingredient").hide();
                    $("#add-ingredient").hide();
                    $("#del-ingredient").hide();
                }
                else {
                    $("#new-ingredient").hide();
                    $("#add-ingredient").hide();
                    $("#del-ingredient").show();
                }
            });

            $("#add-ingredient").on("click", function() {
                var form = {
                    ingredient: $("#new-ingredient").val()
                };
                DrinkMachine.api.create("ingredient", form, function(data) {
                    DrinkMachine.show_alert("success", "Ingredient saved");
                    $("#ingredient").append($("<option>").val(data.ingredient).text(data.ingredient));
                });
            });

            $("#del-ingredient").on("click", function() {
                DrinkMachine.api.remove("ingredient", form.ingredient, function(data) {
                    DrinkMachine.show_alert("success", "Ingredient removed");
                    $("#ingredient option[value='"+data.ingredient+"']").remove();
                });
            });
        };
    };

    this.manage_pumps = new function() {
        this.init = function() {
            $("#pump").on("change", function() {
                if ($(this).val() == 0) return;

                $("#pump-info").show();
                DrinkMachine.api.get("pump", $(this).val(), function(data) {
                    $("#ingredient").val(data.ingredient);
                    $("#flow-rate").val(data.flow_rate);
                });
            });

            $("#save-pump").on("click", function() {
                var form = {
                    id: parseInt($("#pump").val(), 10),
                    ingredient: $("#ingredient").val(),
                    flow_rate: parseFloat($("#flow-rate").val())
                };
                DrinkMachine.api.update("pump", $("#pump").val(), form, function(data) {
                    $("#pump option[value="+data.id+"]").text("Pump " + data.id + " - " + data.ingredient).val(data.id);
                    DrinkMachine.show_alert("success", "Pump info saved");
                });
            });
        };
    };
};
