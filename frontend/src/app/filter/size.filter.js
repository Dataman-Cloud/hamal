export function sizeFilter() {

  return function (rawSize, minUnite) {
    rawSize = parseFloat(rawSize);
    let units = minUnite === "MB" ? ['MB', 'GB', 'TB'] : ['B', 'KB', 'MB', 'GB', 'TB'];
    let unitIndex = 0;
    while (rawSize >= 1024 && unitIndex < units.length - 1) {
      rawSize /= 1024;
      unitIndex++;
    }
    return rawSize.toFixed(2) + ' ' + units[unitIndex];
  }
}
