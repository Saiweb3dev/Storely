export const getSystemInfo = () => {
  return {
    userAgent: navigator.userAgent,
    platform: navigator.mediaDevices,
    language: navigator.language,
    screenResolution: `${window.screen.width}x${window.screen.height}`,
    timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone,
  };
};