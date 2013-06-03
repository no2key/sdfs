CKEDITOR.plugins.add('uploader', {  
    requires: ['dialog'],  
    init: function(a){  
        var b = a.addCommand('uploader', new CKEDITOR.dialogCommand('uploader'));  
        a.ui.addButton('uploader', {  
            label:"uploader",  
            command: 'uploader',  
            icon: this.path + 'icons/uploader.png'  
        });  
        CKEDITOR.dialog.add('uploader', this.path + 'dialogs/uploader.js');  
    }
});  