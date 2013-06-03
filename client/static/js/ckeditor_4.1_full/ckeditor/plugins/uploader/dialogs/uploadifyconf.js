
                  function InsertHTML(v) {
                    // Get the editor instance that we want to interact with.
                    var editor = CKEDITOR.instances.ckeditor;
                    //var value = document.getElementById( 'result' ).value;
                    var value = v;

                    // Check the active editing mode.
                    if ( editor.mode == 'wysiwyg' )
                    {
                      // Insert HTML code.
                      // http://docs.ckeditor.com/#!/api/CKEDITOR.editor-method-insertHtml
                      editor.insertHtml( value );
                    }
                    else
                      alert( '你必须在 WYSIWYG 模式才能在上传之后插入预览内容!' );
                  }

                  $(function() {
                    $('#file_upload').uploadify({
                      'fileObjName':'uploadfile',
                          'debug' :false,
                          'auto':true,
                          'buttonText': "选择上传文件",
                          'removeCompleted':false,
                          'cancelImg': '/static/js/uploadify/uploadify-cancel.png',
                          'swf'      : '/static/js/uploadify/uploadify.swf?ver='+ (new Date()).getTime(),
                          'uploader' : '/root-uploader?ver=' + (new Date()).getTime(),
                          'fileTypeDesc':'支持的格式：',
                  
                          'fileTypeExts':'*.jpg;*.jpge;*.gif;*.png',

                          'overrideEvents' : ['onDialogClose'],
                
                          'onSelect' : function(file) {
                                   
                          },
                          'onSelectError':function(file, errorCode, errorMsg){
                              switch(errorCode) {
                                  case -100:
                                      alert("上传的文件数量已经超出系统限制的"+$('#file_upload').uploadify('settings','queueSizeLimit')+"个文件！");
                                      break;
                                  case -110:
                                      alert("文件 ["+file.name+"] 大小超出系统限制的"+$('#file_upload').uploadify('settings','fileSizeLimit')+"大小！");
                                      break;
                                  case -120:
                                      alert("文件 ["+file.name+"] 大小异常！");
                                      break;
                                  case -130:
                                      alert("文件 ["+file.name+"] 类型不正确！");
                                      break;
                              }
                          },
                          'onFallback':function(){
                              alert("您未安装FLASH控件，无法上传图片！请安装FLASH控件后再试。");
                          },
                          'onUploadSuccess': function(file, data, response) {
                             // InsertHTML("<img src=\"" + data + "\" />");
                              if (data.indexOf('Error') > -1) { 
                                //alert(data);
                                $("#eMessage").append("<div class=\"notification error png_bg\"><div>" + data + "</div></div>");
                              } 
                              else { 
                              //$("#previewImage").attr("src", data.substr(2)).hide().fadeIn(2000); 
                                InsertHTML(data);
                              } 
                            }
                    });
                  });