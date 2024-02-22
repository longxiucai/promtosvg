function uploadFile() {
    var fileInput = document.getElementById('fileInput');
    var file = fileInput.files[0];
    if (!fileInput.files.length) {
      alert('请选择要上传的文件');
      return;
    }

    // 检查文件类型是否为 JSON
    if (!file.type.match('application/json')) {
        alert('请选择 JSON 文件');
        return;
    }

    var formData = new FormData();
    formData.append('file', file);
    var xhr = new XMLHttpRequest();
    var svgPath;
    var iframe = document.getElementById("svgIframe"); // 根据id获取到iframe元素对象
    xhr.open('POST', '/uploadJsonFile', true);
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                console.log('上传成功');
                svgPath = JSON.parse(xhr.response)["svgfile"];
                console.log(svgPath);
                iframe.src = svgPath;
            } else {
                console.log('上传失败');
            }
        }
    };
    xhr.send(formData);
}
