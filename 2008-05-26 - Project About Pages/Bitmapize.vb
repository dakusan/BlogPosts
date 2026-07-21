'Swap in and out bits
Private Declare Function GetDIBits Lib "gdi32" (ByVal aHDC As Long, ByVal hBitmap As Long, ByVal nStartScan As Long, ByVal nNumScans As Long, lpBits As Any, lpBI As BITMAPINFOHEADER, ByVal wUsage As Long) As Long
Private Declare Function SetDIBitsToDevice Lib "gdi32" (ByVal hdc As Long, ByVal x As Long, ByVal y As Long, ByVal dx As Long, ByVal dy As Long, ByVal SrcX As Long, ByVal SrcY As Long, ByVal Scan As Long, ByVal NumScans As Long, Bits As Any, BitsInfo As BITMAPINFOHEADER, ByVal wUsage As Long) As Long
lpBits As Any, lpBitsInfo As BITMAPINFOHEADER, ByVal wUsage As Long, ByVal dwRop As Long) As Long
Private Type RGBQUAD
		b As Byte
		g As Byte
		r As Byte
		Reserved As Byte
End Type
Private Type BITMAPINFOHEADER '40 bytes
		biSize As Long
		biWidth As Long
		biHeight As Long
		biPlanes As Integer
		biBitCount As Integer
		biCompression As Long
		biSizeImage As Long
		biXPelsPerMeter As Long
		biYPelsPerMeter As Long
		biClrUsed As Long
		biClrImportant As Long
End Type
Private Const DIB_RGB_COLORS = 0 '  color table in RGBs

'Prepare colors
Private Declare Sub CopyMemory Lib "kernel32" Alias "RtlMoveMemory" (Destination As Any, Source As Any, ByVal Length As Long)
Private Declare Function GetBkColor Lib "gdi32" (ByVal hdc As Long) As Long

Public Sub DisplaySignature(ByRef TheForm As Form)
    'Read in Signature
    Dim BitmapLength As Long, OutBitmap() As RGBQUAD, BitInfo As BITMAPINFOHEADER, Signature As PictureBox
    Set Signature = TheForm.Signature
    BitmapLength = Signature.Width * Signature.Height
    ReDim OutBitmap(0 To BitmapLength - 1) As RGBQUAD
    With BitInfo
            .biSize = 40
            .biWidth = Signature.Width
            .biHeight = -Signature.Height
            .biPlanes = 1
            .biBitCount = 32
            .biCompression = 0 'BI_RGB
            .biSizeImage = .biWidth * 4 * -.biHeight
    End With
    GetDIBits Signature.hdc, Signature.Image, 0, Signature.Height, OutBitmap(0), BitInfo, DIB_RGB_COLORS
    
    'Alpha blend signature
    Dim i As Long, Alpha As Double, BackColor As RGBQUAD, ForeColor As RGBQUAD, OBC As Long, OFC As Long
    OFC = &H67906
    OBC = GetBkColor(TheForm.hdc)
    CopyMemory BackColor, OBC, 4
    CopyMemory ForeColor, OFC, 4
    For i = 0 To BitmapLength - 1
        Alpha = 1 - (CDbl(OutBitmap(i).r) / 255)
        OutBitmap(i).r = ForeColor.r * Alpha + BackColor.r * (1 - Alpha)
        OutBitmap(i).g = ForeColor.g * Alpha + BackColor.g * (1 - Alpha)
        OutBitmap(i).b = ForeColor.b * Alpha + BackColor.b * (1 - Alpha)
    Next i
    
    SetDIBitsToDevice TheForm.hdc, Signature.Left, Signature.Top, Signature.Width, Signature.Height, 0, 0, 0, Signature.Height, OutBitmap(0), BitInfo, DIB_RGB_COLORS
    TheForm.Refresh
End Sub