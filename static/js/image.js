function insertImageToCmtContent(imagePath) {
    //直接追加一个图片标记，style需要用户自行添加
    const element = document.getElementById("EditingInputControl");
    element.value = element.value + "[img]" + " src=\"" + imagePath + "\" style=\"\" " + "[/img]";
}

function requestImagesAndShow(pageIndex) {
    const httpRequest = new XMLHttpRequest();
    httpRequest.open("GET","/img?ImagePageIndex="+pageIndex,true);
    httpRequest.onreadystatechange=function(){
        //没完成，就不做任何操作
        if(httpRequest.readyState !== 4){
            return;
        }
        document.getElementById("ImagesForSelectToInsert").innerHTML = httpRequest.response
    };
    httpRequest.send(null);
}