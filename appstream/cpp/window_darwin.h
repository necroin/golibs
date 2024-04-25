#include <libproc.h>
#include <sys/proc_info.h>
#include <CoreGraphics/CoreGraphics.h>

uint32_t GetActiveDisplaysCount()
{
    CGDirectDisplayID displays[1024];
    uint32_t count = 0;
    CGGetActiveDisplayList(1024, displays, &count);
    return count;
}

CGDirectDisplayID *GetActiveDisplays()
{
    CGDirectDisplayID displays[1024];
    uint32_t count = 0;
    CGGetActiveDisplayList(1024, displays, &count);
    return displays;
}

CGDirectDisplayID IndexArray(CGDirectDisplayID *displays, uint32_t index)
{
    return displays[index];
}

typedef struct WindowInfo
{
    int pid;
    char *windowName;
    CGRect bounds;

} WindowInfo;

CFIndex GetWindowsCount()
{
    CFArrayRef windows = CGWindowListCopyWindowInfo(kCGWindowListExcludeDesktopElements, kCGNullWindowID);
    CFIndex windowsCount = CFArrayGetCount(windows);
    return windowsCount;
}

WindowInfo *GetWindowsInfo()
{
    CFArrayRef windows = CGWindowListCopyWindowInfo(kCGWindowListOptionAll, kCGNullWindowID);
    CFIndex windowsCount = CFArrayGetCount(windows);
    WindowInfo *infos = (WindowInfo *)calloc(windowsCount, sizeof(WindowInfo));
    for (size_t i = 0; i < windowsCount; i++)
    {
        CFDictionaryRef windowInfo = (CFDictionaryRef)CFArrayGetValueAtIndex(windows, i);

        int windowPid;
        CFNumberRef windowPidRef = (CFNumberRef)CFDictionaryGetValue(windowInfo, kCGWindowOwnerPID);
        CFNumberGetValue(windowPidRef, kCFNumberIntType, &windowPid);

        char windowName[256];
        windowName[0] = 0;
        CFStringRef windowNameRef = (CFStringRef)CFDictionaryGetValue(windowInfo, kCGWindowName);
        CFStringGetCString(windowNameRef, windowName, 256, kCFStringEncodingUTF8);

        CGRect bounds;
        CFDictionaryRef boundsRef = (CFDictionaryRef)CFDictionaryGetValue(windowInfo, kCGWindowBounds);
        CGRectMakeWithDictionaryRepresentation(boundsRef, &bounds);

        WindowInfo info;
        info.pid = windowPid;
        info.windowName = windowName;
        info.bounds = bounds;
        infos[i] = info;
    }
    return infos;
}

WindowInfo IndexWindowInfo(WindowInfo *windows, CFIndex index)
{
    return windows[index];
}
