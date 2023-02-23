Number.prototype.formatBytes = function () {
    var units = ['B', 'KB', 'MB', 'GB', 'TB'],
        bytes = this,
        i;
    for (i = 0; bytes >= 1024 && i < 4; i++) {
        bytes /= 1024;
    }
    return bytes.toFixed(2) + units[i];
}
Number.prototype.secondsToHms = function () {
    var d = Math.floor(this / (3600 * 24));
    var h = Math.floor(this % (3600 * 24) / 3600);
    var m = Math.floor(this % 3600 / 60);
    var s = Math.floor(this % 60);

    return (d > 0 ? d + 'd ' : "")
        + (h > 0 ? h + 'h ' : "")
        + (m > 0 ? m + 'm ' : "")
        + (s > 0 ? s + 's ' : "");
}

document.addEventListener('alpine:init', () => {
    Alpine.data('systeminfo', () => ({
        data: {
            "percent": { "cpu": 0.0, "mem": 0.0, "swap": 0.0, "disk": 0.0 },
            "load": { "load1": 0, "load5": 0, "load15": 0 },
            "memory": { "total": 0, "used": 0, "free": 0 },
            "swap": { "total": 0, "used": 0, "free": 0 },
            "disk": { "total": 0, "used": 0, "free": 0, "read": 0, "write": 0 },
            "network": { "rx": 0, "tx": 0, "in": 0, "out": 0, "min": 0, "mout": 0 },
            "uptime": 0,
            "hostname": ''
        },
        updatetime: new Date().toLocaleString(),
        thresholdColor() {
            const thresholds = [
                { val: 80, color: 'red' },
                { val: 60, color: 'yellow' },
                { val: 0, color: 'green' },
            ]
            for (var i = 0; i < thresholds.length; i++) {
                for (var name in this.data.percent) {
                    if (this.data.percent[name] > thresholds[i].val) {
                        return thresholds[i].color;
                    }
                }
            }
            return 'black';
        },
        init() {
            this.load();
            window.setInterval(() => this.load(), 3000);
        },
        load() {
            fetch("/stat").then(response => response.json())
                .then(data => {
                    this.data = data;
                    this.updatetime = new Date().toLocaleString();
                });
        }
    }))
})