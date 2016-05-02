@echo off
setlocal
if not "%VS120COMNTOOLS%" == "" (
  set "VSCOMNTOOLS=%VS120COMNTOOLS%"
) else if not "%VS140COMNTOOLS%" == "" (
  set "VSCOMNTOOLS=%VS140COMNTOOLS%"
)

call "%VSCOMNTOOLS%..\..\VC\bin\vcvars32.bat"
cl.exe hello.c /Z7 /Fohello32.obj

call "%VSCOMNTOOLS%..\..\VC\bin\amd64\vcvars64.bat"
cl.exe hello.c /Z7 /Fohello64.obj
