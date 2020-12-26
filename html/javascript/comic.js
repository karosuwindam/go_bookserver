filetype="";
var lefte =[0];
var mide = [1,3,4];
var righte = [2];
var pagemax = 4;            //pagemax
var pagenow = 0;
var modeonetwe = true;      //2 page view
var filepass = ""      //ファイルパス
var filename = "";          //ファイルネーム
var filetype = ".jpg";      //ファイル拡張子
var mouseover_f = true;     //mouns over flag
var loaddata = true;
var view_message = ["Next","1 page next","back","1/2 page change","page view"]
document.onkeydown = function(event){
    var keyEvent = event||window.event;
    console.log("KEYEVENT:"+keyEvent.keyCode);
    switch (keyEvent.keyCode){
        case 37:			//←
            nextpage(+2);
            break;
        case 38:			//↑
            nextpage(+1);
            break;
        case 39:			//→
            nextpage(-2);
            break;
        case 40:			//↓
            nextpage(-1);
            break;
        case 72:            //h
            helpviewonoff();
            break;
        case 73:        //i
            infoviewonoff();
            break;
        case 77:        //m
            mouseover_f = !mouseover_f;
            break;
        case 80:        //p
            pageviewonoff()
            break;
        case 32:       //space
            onclickevent(3);
            break;
        default:
            break;
    }
}

function borderviewonoff(){
    var flag = false
    var data = ["leftT","midT","rightT","midM","midB"];
    for(var i=0;i<data.length;i++){
        var tmp = document.getElementsByClassName(data[i]);
        if (i==0&&(tmp[0].style.border =="")){
            flag = true
        }
        if(flag){
            tmp[0].style.border = "1px solid  #000000";
        }else{
            tmp[0].style.border = "";
        }
    }//border: inset 10px #ff0000;

}
function infoviewonoff(){
    var info = document.getElementsByClassName("output");
    if (info[0].style.display == ""){
        info[0].style.display = "none"
    }else{
        info[0].style.display = ""
    }
}
function helpviewonoff(){
    var help = document.getElementsByClassName("help");
    if (help[0].style.display == ""){
        help[0].style.display = "none"
    }else{
        help[0].style.display = ""
    }
}
function pageviewonoff(){
    var info = document.getElementsByClassName("nowpage");
    var sidber = document.getElementById("pageslider");
    if (info[0].style.display == ""){
        info[0].style.display = "none"
        sidber.style.display = "none"
    }else{
        info[0].style.display = ""
        sidber.style.display = ""
    }
}
function mouseover(num,class_t)
{
    var tmp_data = ""
    if (mouseover_f){
        switch(num){
            case 0:     //left
            case 1: //midT
            case 2:     //right
                document.getElementsByClassName(class_t)[0].innerHTML = view_message[num]
                break;
            case 3: //midM
                if(modeonetwe){
                    tmp_data = "1 page view"
                }else{
                    tmp_data = "2 page view"
                }
                document.getElementsByClassName(class_t)[0].innerHTML = tmp_data//"1 to 2 change"
                break;
            case 4: //midB
                document.getElementsByClassName(class_t)[0].innerHTML = view_message[num]
                break;
            default:
                break;
        }
    }
}
function mouseout(num,class_t)
{
    switch(num){
        case 0:     //left
        case 1: //midT
        case 2:     //right
        case 3: //midM
        case 4: //midB
            document.getElementsByClassName(class_t)[0].innerHTML = ""
            break;
        default:
            break;
    }
}
function onclickevent(num)
{
    switch(num){
        case 0:
            nextpage(2)
            break;
        case 1:
            nextpage(1)
            break;
        case 2:
            nextpage(-2)
            break;
        case 3:
            chPageOneTwe(2);
            break;
        case 4:
            pageviewonoff();
            // chIndexView();
            break;
        default:
            break;
    }
    //document.getElementById("output").innerHTML = num+" page:"+pagenow;
}
var onloadclose = function(){     //起動時を非表示項目
    var info = document.getElementsByClassName("help");
    var page = document.getElementsByClassName("nowpage");
    var sidber = document.getElementById("pageslider");
    info[0].style.display = "none"
    page[0].style.display = "none"
    sidber.style.display = "none"

    var data = ["leftT","midT","rightT","midM","midB"];
    for(var i=0;i<data.length;i++){
        var tmp = document.getElementsByClassName(data[i]);
        tmp[0].style.border = "";
        tmp[0].innerHTML = "";
    }
}
function mousescroll_int(){
    //Firefox
    if(window.addEventListener){
        window.addEventListener('DOMMouseScroll', function(e){
            // alert(e.detail);
            if((document.getElementsByClassName("listdata")[0].style.display == "none")
            &&(document.getElementsByClassName("output")[0].style.display == "none")){
                if (e.detail < 0){
                    nextpage(-2);
                }else{
                    nextpage(+2);
                }
            }
        }, false);
    }
    
    //IE
    if(document.attachEvent){
        document.attachEvent('onmousewheel', function(e){
            // alert(e.wheelDelta);
            if((document.getElementsByClassName("listdata")[0].style.display == "none")
            &&(document.getElementsByClassName("output")[0].style.display == "none")){
                if (e.wheelDelta > 0){
                    nextpage(-2);
                }else{
                    nextpage(+2);
                }
            }
        });
    }
    
    //Chrome
    window.onmousewheel = function(e){
        // alert(e.wheelDelta);
        if((document.getElementsByClassName("listdata")[0].style.display == "none")
        &&(document.getElementsByClassName("output")[0].style.display == "none")){
            if (e.wheelDelta > 0){
                nextpage(-2);
            }else{
                nextpage(+2);
            }
        }
    }

}
function initload(){
    var data = ["leftT","midT","rightT","midM","midB"];
    for(var i=0;i<data.length;i++){
        var tmp = document.getElementsByClassName(data[i]);
        tmp[0].style.border = "5px solid  #000000";
        tmp[0].innerHTML = view_message[i];
    }//border: inset 10px #ff0000;
    mousescroll_int()
    nowpage(pagenow);
}
window.onload=function() {          //起動後
    this.initload();
    setTimeout(onloadclose, 3000);
};
function nowpage(num){
    var r_imgtmp = new Array(3);
    var l_imgtmp = new Array(3);
    for (var i=0;i<r_imgtmp.length;i++){
        r_imgtmp[i] = document.getElementById("pagei0_"+i);
        l_imgtmp[i] = document.getElementById("pagei1_"+i);
    }
    pagenow = num-0;
    if (pagenow==0){
        var r_page =[pagemax-1,0,1]
        var l_page = [pagemax,0,2]
        for(var i=0;i<r_page.length;i++){
            if (i==1){
                r_imgtmp[i].src = "/img/blank.jpg"
            }else{
                r_imgtmp[i].src = filepass + filename + r_page[i] + filetype;
            }
            l_imgtmp[i].src = filepass + filename + l_page[i] + filetype;
        }
        return
    }
    for(var i=0;i<r_imgtmp.length;i++){
        var r_page_tmp = pagenow +(i-1)*2;
        var l_page_tmp = pagenow +(i-1)*2+1;
        var flag = [false,false]
        if (r_page_tmp < 0){
            r_page_tmp = pagemax-1;
        }else if(r_page_tmp > pagemax){
            r_page_tmp = 0;
            flag[0] = true;
        }
        if (l_page_tmp < 0){
            l_page_tmp = pagemax;
        }else if(l_page_tmp > pagemax){
            l_page_tmp = 1;
            flag[1] = true;
        }
        if(loaddata){
            if(i==1){
                flag[0]=true;
                l_page_tmp = 0;
            }
        }

        if(flag[0]&&i==1){
            r_imgtmp[i].src = "/img/blank.jpg"
        }else{
            r_imgtmp[i].src = filepass + filename + r_page_tmp + filetype;
        }
        if(flag[1]&&i==1){
            l_imgtmp[i].src = "/img/blank.jpg"
        }else{
            l_imgtmp[i].src = filepass + filename + l_page_tmp + filetype;
        }

    }
    pagerenew();

}
function nextpage(num){
    var r_imgtmp = new Array(3);
    var l_imgtmp = new Array(3);
    for (var i=0;i<r_imgtmp.length;i++){
        r_imgtmp[i] = document.getElementById("pagei0_"+i);
        l_imgtmp[i] = document.getElementById("pagei1_"+i);
    }
    if(((pagenow == 0)&&(loaddata))||(!modeonetwe) )
    {
        if (num >0){
            pagenow++;
        }else{
            pagenow--;
        }
    }else{
        pagenow += num;
    }
    if (pagenow>pagemax)
    {
        pagenow=0;
    }else if(pagenow<0){
        pagenow=pagemax;
    }
    for(var i=0;i<r_imgtmp.length;i++){
        var r_page_tmp = pagenow +(i-1)*2;
        var l_page_tmp = pagenow +(i-1)*2+1;
        var flag = [false,false]
        if (r_page_tmp < 0){
            r_page_tmp = pagemax;
        }else if(r_page_tmp > pagemax){
            r_page_tmp = 0;
            flag[0] = true;
        }
        if (l_page_tmp < 0){
            l_page_tmp = pagemax;
        }else if(l_page_tmp > pagemax){
            l_page_tmp = 1;
            flag[1] = true;
        }

        if(flag[0]&&i==1){
            r_imgtmp[i].src = "/img/blank.jpg"
        }else{
            r_imgtmp[i].src = filepass + filename + r_page_tmp + filetype;
        }
        if(flag[1]&&i==1){
            l_imgtmp[i].src = "/img/blank.jpg"
        }else{
            l_imgtmp[i].src = filepass + filename + l_page_tmp + filetype;
        }

    }
    document.getElementById("pageslider").value = pagenow;
    loaddata = false;
    pagerenew();
}
function pagerenew(){
    var output = document.getElementsByClassName("nowpage");
    output[0].innerHTML = (pagenow+1)+"/"+(pagemax+1) +" page"
}
function chPageOneTwe(mode){
    var r_imgtmp = new Array(3);
    var l_imgtmp = new Array(3);
    for (var i=0;i<r_imgtmp.length;i++){
        r_imgtmp[i] = document.getElementById("page0_"+i);
        l_imgtmp[i] = document.getElementById("page1_"+i);
    }
    if(mode == 0){
        modeonetwe = false;
    }else if(mode == 1){
        modeonetwe = true;
    }else{
        modeonetwe = !modeonetwe;
    }
    if (modeonetwe){        //切り替え有効 2page
        for (var i = 0;i < r_imgtmp.length ;i++){
            r_imgtmp[i].style.left = "50%";
            r_imgtmp[i].style.width = "45%";
            r_imgtmp[i].style.textAlign = "left";
        }
        for (var i = 0;i < l_imgtmp.length ;i++){
            if(i==1)
            {
                l_imgtmp[i].style.display = ""
            }
        }
        
    }else{      //1page 表示
        for (var i = 0;i < r_imgtmp.length ;i++){
            r_imgtmp[i].style.left = "5%";
            r_imgtmp[i].style.width = "90%";
            r_imgtmp[i].style.textAlign = "center";
        }
        for (var i = 0;i < l_imgtmp.length ;i++){
            l_imgtmp[i].style.display = "none"
        }
    }
    if (pagenow == 0){
        if(modeonetwe){
            document.getElementById("pagei0_1").src = "/img/blank.jpg";
            document.getElementById("pagei1_1").src = filepass + filename + (pagenow) + filetype;
        }else{
            document.getElementById("pagei0_1").src = filepass + filename + pagenow + filetype;
            document.getElementById("pagei1_1").src = filepass + filename + (pagenow+1) + filetype;
        }
    }
}

function chIndexView(){
    var tmp = document.getElementsByClassName("output");
    if (tmp[0].style.display == ""){
        tmp[0].style.display = "none"
    }else{
        tmp[0].style.display = ""
    }
}


/*
 * スワイプイベント設定
 */

function setSwipe() {
	let startX;		// タッチ開始 x座標
	let startY;		// タッチ開始 y座標
	let moveX;	// スワイプ中の x座標
	let moveY;	// スワイプ中の y座標
	let dist = 30;	// スワイプを感知する最低距離（ピクセル単位）
	
	// タッチ開始時： xy座標を取得
	window.addEventListener("touchstart", function(e) {
		e.preventDefault();
		startX = e.touches[0].pageX;
		startY = e.touches[0].pageY;
	});
	
	// スワイプ中： xy座標を取得
	window.addEventListener("touchmove", function(e) {
		e.preventDefault();
		moveX = e.changedTouches[0].pageX;
		moveY = e.changedTouches[0].pageY;
	});
	
	// タッチ終了時： スワイプした距離から左右どちらにスワイプしたかを判定する/距離が短い場合何もしない
	window.addEventListener("touchend", function(e) {
        if (startX > moveX && startX > moveX + dist) {		// 右から左にスワイプ
            // console.log("left");
            nextpage(-2);
		}
		else if (startX < moveX && startX + dist < moveX) {	// 左から右にスワイプ
            // console.log("right");
            nextpage(+2);
		}
	});
}
window.addEventListener("load", function(){

	// スワイプイベント設定
	setSwipe();
});
