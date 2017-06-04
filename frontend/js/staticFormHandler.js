(function($){
  $.fn.staticFormHandler = function(options){

    var defaults = {
        name: 'staticFormHandler',
        fields:{},
        url:"",
        confirmation_txt:"Thanks!",
        error_txt:"Sorry, couldn't connect to the server. If you try again later, it might work.",
        display_condition:function(){ return false; },
        container:"#formContainer",
        ask_only_once: true
      },
      plugin = this;

    this.opt = function(field){
      return options[field] || defaults[field];
    };

    this.createForm = function(){
      var $f = $("<form action="+ plugin.opt("action") + ">")
      $f.append(plugin.getFormFields());
      return $f;
    };

    this.readCookie = function(){
      c = document.cookie.split('; ');
      var i = c.length-1;
      while (i --> 0){
        var C = c[i].split('=');
        if(C[0] === plugin.opt('name')){
          return C[1];
        }
      }
    };

    this.getCookieExpireDate = function(){
      var now = new Date();
      var expDate = new Date();
      expDate.setYear(now.getFullYear()+20);
      return expDate.toString();
    };

    this.setCookie = function(){
      //document.cookie = plugin.opt('name')+"=seen;expires=" + plugin.getCookieExpireDate();
    };

    this.alreadySeen = function(){
      return this.readCookie() === 'seen';
    };

    this.getInputField = function(f, c){
       return '<div><label for="'+f+'"></label><input type="text" name="'+f+'" value="" id="'+f+'"></div>';
    };

    this.getSendButton = function(){
       var $btn = $('<a href="." style="border: 1px solid black; padding:5px; margin-top:5px;"> send </a>');
      $btn.click(plugin.sendData);
       return $btn;
    };

    this.getFormFields = function(){
      var fields = plugin.opt("fields"),
          fields_html = "";
      for( var f in fields){
        switch(fields[f].type){
          case "input":
            fields_html += plugin.getInputField(f, fields[f]);
            break;
        }
      }

      return $(fields_html).append(plugin.getSendButton());
    };

    this.clearInterval = function(){
      if( typeof(plugin.ti) !== 'undefined'){
        clearInterval(plugin.ti);
      }
    };

    this.getTriggerFunction = function(){
      return function(){
        if(plugin.opt('display_condition')()) {
          plugin.showForm();
          plugin.setCookie();
          plugin.clearInterval();
        }
      };
    };

    this.gatherData = function() {
      var fields = plugin.opt("fields"),
          data = {};
      for( var f in fields){
        data[f] = $(f).val();
      }
      return data;
    };

    this.onAjaxError = function (err){
      console.log(err);
    };

    this.onAjaxSuccess = function (data) {
      console.log(data);
    };

    this.sendData = function (){
      console.log('seding data to');
      console.log(plugin.opt('url'))
       $.ajax({
        method: "PUT",
        url: plugin.opt('url'),
        data: '{"Email":"ingmar@drewing.de"}',
        dataType: "json",
         contentType: "application/json",
        error: plugin.onAjaxError,
        success: plugin.onAjaxSuccess
      });
      return false;
    };

    this.showForm = function(){ 
      this.each(function(){
        $(this).append( plugin.createForm());
      });
    };

    if( ! plugin.alreadySeen() ){
      plugin.ti = setInterval(plugin.getTriggerFunction(), 100);
    }
  };
})(jQuery)
