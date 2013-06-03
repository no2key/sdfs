CKEDITOR.dialog.add('uploader', function(editor){
    var escape = function(value){
        return value;
    };
    return {
        title: '上传图片',
        resizable: CKEDITOR.DIALOG_RESIZE_BOTH,
        minWidth: 720,
        minHeight: 480,
        contents: [{
            id: 'cb',
            name: 'cb',
            label: 'cb',
            title: 'cb',
            elements: [{
            		id: 'imger', 
					type : 'html',
					html : '<div><form> <div id="queue"></div>  <input id="file_upload" name="file_upload" type="file" multiple="true"></form><script type="text/javascript">$(function(){$(\'#file_upload\').uploadify({\'fileObjName\':\'uploadfile\',\'debug\':false,\'auto\':true,\'buttonText\':"选择上传文件",\'removeCompleted\':false,\'cancelImg\':\'/static/js/uploadify/uploadify-cancel.png\',\'swf\':\'/static/js/uploadify/uploadify.swf?ver=\'+(new Date()).getTime(),\'uploader\':\'/root-uploader?ver=\'+(new Date()).getTime(),\'fileTypeDesc\':\'支持的格式：\',\'fileTypeExts\':\'*.jpg;*.jpge;*.gif;*.png\',\'overrideEvents\':[\'onDialogClose\'],\'onSelect\':function(a){},\'onSelectError\':function(a,b,c){switch(b){case-100:alert("上传的文件数量已经超出系统限制的"+$(\'#file_upload\').uploadify(\'settings\',\'queueSizeLimit\')+"个文件！");break;case-110:alert("文件 ["+a.name+"] 大小超出系统限制的"+$(\'#file_upload\').uploadify(\'settings\',\'fileSizeLimit\')+"大小！");break;case-120:alert("文件 ["+a.name+"] 大小异常！");break;case-130:alert("文件 ["+a.name+"] 类型不正确！");break}},\'onFallback\':function(){alert("您未安装FLASH控件，无法上传图片！请安装FLASH控件后再试。")},\'onUploadSuccess\':function(a,b,c){if(b.indexOf(\'Error\')>-1){$("#result").append("<h1>"+b+"</h1>")}else{$("#result").append("<div>"+b+"</div>")}}})});</script></div>'
				}]
        }],
            onHide: function () {
                alert('onHide');
            },
	        onOk: function(){

	            editor.insertHtml("<pre class=\"brush:;\">" + "html" + "</pre>");
	        },
            onCancel: function () {
                alert('onCancel');
            }
    };
});