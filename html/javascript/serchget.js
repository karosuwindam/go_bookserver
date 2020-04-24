var serchurl = ["bookname","filelist","copyfile"]
var selectdata = 0

function bookname_edit(ary){
    var output =[]
    for (var i=0;i<ary.length;i++){
        var tmp =[]
        tmp.push(ary[i].id,ary[i].name,ary[i].title,ary[i].Writer,ary[i].brand,ary[i].ext)
        output.push(tmp)
    }
    return output
}
function copyfile_edit(ary){
    var output =[]
    for (var i=0;i<ary.length;i++){
        var tmp =[]
        tmp.push(ary[i].id,ary[i].Zippass,ary[i].Filesize,ary[i].copyflag)
        output.push(tmp)
    }
    return output
}
function filelist_edit(ary){
    var output =[]
    for (var i=0;i<ary.length;i++){
        var tmp =[]
        tmp.push(ary[i].id,ary[i].name,ary[i].pdfpass,ary[i].Zippass,ary[i].tag)
        output.push(tmp)
    }
    return output
}
function outputTable(ary){
    var output
    var sumdata = 0
    var outf =false
    output = "<table>"
    for(var i=0;i<ary.length;i++){
        output += "<tr style='color:#FFFFFF'>";
        var tmp = ary[i]
        for(var j=0;j<tmp.length;j++){
            output += "<td>"+tmp[j]+"</td>";
            outf = true;
        }
        if(selectdata == 2){
            sumdata += tmp[2]-0
            output += "<td><input type='checkbox' "
            if (tmp[3] == "1"){
              output += "checked='checked' "
            }
            output += "onclick=\"ck_copyfilebox(\'"+tmp[1]+"\',this.checked)\""
            output += ">"+"</td>"
        }
        output += "</tr>"
    }
    output += "</table>"
    if((selectdata == 2)){
        output += "<br>"
        if (sumdata > 1024*1024*1024){
            output += "<div>"+"sumsize:"+Math.round(sumdata/1024/1024/1024*1000)/1000 +"G"+"</div>"
        }else if (sumdata > 1024*1024){
            output += "<div>"+"sumsize:"+Math.round(sumdata/1024/1024*1000)/1000 +"M"+"</div>"
        }else if (sumdata > 1024){
            output += "<div>"+"sumsize:"+Math.round(sumdata/1024*1000)/1000 +"K"+"</div>"
        }else{
            output += "<div>"+"sumsize:"+sumdata+"</div>"
        }
    }
    if (!outf){
        output = ""
    }
    return output
}
function outputhtmlJson(str){
    var tmp = JSON.parse(str);
    var ary;
    console.log(tmp)
    switch (selectdata){
        case 0:
            ary = bookname_edit(tmp);
            break;
        case 1:
            ary = filelist_edit(tmp);
            break;
        case 2:
            ary = copyfile_edit(tmp);
            break;
    }
    console.log(ary)
    return outputTable(ary)
}

function serchgetData(output){
    var out = document.getElementById(output)
    var url
    url = "/serch/" + serchurl[selectdata] +"/"
    
    var keyword = document.getElementById("serch").value
    var req = new XMLHttpRequest();
    req.onreadystatechange = function(){
        if(req.readyState == 4 && req.status == 200){
            nowserchpage = 1
            var data=req.responseText;
            out.innerHTML = outputhtmlJson(data)
            //out.innerHTML = data
        }
    };
    req.open("GET",url+keyword,false);
    req.send(null);
}