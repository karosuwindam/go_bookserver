var meta_suburl=""
function destory(id){
  myRet = confirm("destory id="+id+" OK??");
  if (myRet){
    var xhr = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    xhr.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(xhr.readyState == 4 && xhr.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = xhr.responseText;
          console.log(data);		          // 取得した ファイルの中身を表示
          document.getElementById("answer").innerHTML = "destory id=" + id + " OK"
          getJSON("/list","output");
      }
    };

    var url = "destory/"
    if (meta_suburl != ""){
      url += meta_suburl + "/"
    }
    url += id;
    xhr.open('POST', url, true);
    xhr.setRequestHeader('content-type', 'application/x-www-form-urlencoded;charset=UTF-8');
    xhr.send(null);
  }
}

function jsonOutput(str){
    var output = ""
    var table_title = ["id","name","title"]
    if (meta_suburl=="filelist"){
      table_title = ["id","name","pdfpass","zippass","tag"]
    }
    var tmp = JSON.parse(str)
    output += "<table>"
    output += "<tr>"
    for (var i=0;i<table_title.length;i++){
        output += "<th>"+table_title[i]+"</th>"
    }
    output += "</tr>"
    for (var i=0; i< tmp.length;i++){
    //   output += "<div>"
      output += "<tr>"
    //   output += tmp[i].name
      output += "<td>"+tmp[i].id+"</td>"
      output += "<td>"+tmp[i].name+"</td>"
      if (meta_suburl=="filelist"){
        output += "<td>"+tmp[i].pdfpass+"</td>"
        output += "<td>"+tmp[i].Zippass+"</td>"
        output += "<td>"+tmp[i].tag+"</td>"
      }else{
        output += "<td>"+tmp[i].title+"</td>"
      }
    //   output += " <a href='edit/"+tmp[i].id+"'>"+"edit"+"</a>"
      output += "<td><a href='show/"
      if (meta_suburl!=""){
        output +=meta_suburl+"/"
      }
      output += tmp[i].id+"'>"+"show"+"</a></td>"
      output += "<td><a href='edit/"
      if (meta_suburl!=""){
        output +=meta_suburl+"/"
      }
      output += tmp[i].id+"'>"+"edit"+"</a></td>"
    //   output += " <a href='destory/"+tmp[i].id+"'>"+"destory"+"</a>"
      output += "<td><a href='javascript:destory("+tmp[i].id+");'>"+"destory"+"</a></td>"
      output += "</tr>"
    //   output +="</div>"
    }
    output += "</table>"
    return output
  }
  function getJSON(url,output) {
    var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = req.responseText;
          console.log(data);		          // 取得した JSON ファイルの中身を表示
          document.getElementById(output).innerHTML = jsonOutput(data);
      }
    };
    req.open("GET", url+"/"+meta_suburl, false); // HTTPメソッドとアクセスするサーバーの　URL　を指定
    req.send(null);					    // 実際にサーバーへリクエストを送信
  }

  function serchgetJSON(output){
    var keyword = document.getElementById("keyword").value
    var req = new XMLHttpRequest();
    req.onreadystatechange = function(){
      if(req.readyState == 4 && req.status == 200){
        var data=req.responseText;
        document.getElementById(output).innerHTML = data;
      }
    };
    req.open("GET","/serch/filelist/"+keyword,false);
    req.send(null);
  }