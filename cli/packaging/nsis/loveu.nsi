; loveu CLI Windows installer (NSIS)
; Build: makensis /DVERSION=0.1.0 /DOUTFILE=loveu-0.1.0-windows-amd64-setup.exe /DSOURCE=path\to\loveu.exe loveu.nsi

!ifndef VERSION
  !define VERSION "0.1.0"
!endif
!ifndef OUTFILE
  !define OUTFILE "loveu-${VERSION}-windows-amd64-setup.exe"
!endif
!ifndef SOURCE
  !define SOURCE "loveu.exe"
!endif

Name "loveu ${VERSION}"
OutFile "${OUTFILE}"
InstallDir "$PROGRAMFILES64\loveu"
RequestExecutionLevel admin
Unicode true

!include "MUI2.nsh"
!include "WinMessages.nsh"

!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_LANGUAGE "English"

!macro AppendToPath DIR
  ReadRegStr $0 HKLM "SYSTEM\CurrentControlSet\Control\Session Manager\Environment" "Path"
  StrCmp $0 "" AppendEmpty
  ; skip if already present
  Push "$0"
  Push "${DIR}"
  Call StrContains
  Pop $1
  StrCmp $1 "" 0 AppendDone
  WriteRegExpandStr HKLM "SYSTEM\CurrentControlSet\Control\Session Manager\Environment" "Path" "$0;${DIR}"
  Goto Broadcast
AppendEmpty:
  WriteRegExpandStr HKLM "SYSTEM\CurrentControlSet\Control\Session Manager\Environment" "Path" "${DIR}"
Broadcast:
  SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000
AppendDone:
!macroend

Function StrContains
  Exch $R1 ; needle
  Exch
  Exch $R2 ; haystack
  Push $R3
  Push $R4
  Push $R5
  StrLen $R3 $R1
  StrCpy $R4 0
loop:
  StrCpy $R5 $R2 $R3 $R4
  StrCmp $R5 $R1 found
  StrCmp $R5 "" notfound
  IntOp $R4 $R4 + 1
  Goto loop
found:
  StrCpy $R1 $R1
  Goto done
notfound:
  StrCpy $R1 ""
done:
  Pop $R5
  Pop $R4
  Pop $R3
  Pop $R2
  Exch $R1
FunctionEnd

Section "Install"
  SetOutPath "$INSTDIR"
  File /oname=loveu.exe "${SOURCE}"
  WriteUninstaller "$INSTDIR\Uninstall.exe"
  !insertmacro AppendToPath "$INSTDIR"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\loveu" "DisplayName" "loveu ${VERSION}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\loveu" "UninstallString" "$\"$INSTDIR\Uninstall.exe$\""
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\loveu" "DisplayVersion" "${VERSION}"
  WriteRegStr HKLM "Software\loveu" "InstallDir" "$INSTDIR"
SectionEnd

Section "Uninstall"
  Delete "$INSTDIR\loveu.exe"
  Delete "$INSTDIR\Uninstall.exe"
  RMDir "$INSTDIR"
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\loveu"
  DeleteRegKey HKLM "Software\loveu"
SectionEnd
