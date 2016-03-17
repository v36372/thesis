package main

import (
    "fmt"
    "github.com/mjibson/go-dsp/wav"
    "github.com/mjibson/go-dsp/fft"
    wavwriter "github.com/cryptix/wav"
    "os"
    "math"
    "math/cmplx"
    "strconv"
)

const (
  FREQ_THRES = 0.0001
  BIT_OFFSET = 1
  SAMPLE_PER_FRAME=3000
  BIT_REPEAT=5
)

func main() {
    file , _ := os.Open("test.wav")
    reader , _ := wav.New(file)

    l,err := reader.ReadFloats(reader.Samples)
    if err!= nil {
      fmt.Println(err)
    }
    fmt.Println(reader.Header.ByteRate)
    fmt.Println(reader.Header.SampleRate)
    //---------------------fft the whole file-------------------------
    mag := make([]float64, reader.Samples)
    phs := make([]float64, reader.Samples)

    fourier := fft.FFTReal(l)
    for i,a:= range fourier{
        mag[i], phs[i] = cmplx.Polar(a)
    }
    //----------------------fft the whole file------------------------

    //==================================================================

    //-------------------divide into frames--------------------------
    // mag := make([]float64, 0)
    // phs := make([]float64, 0)
    //
    // var max float64
    // max = 0
    // var i = SAMPLE_PER_FRAME-1
    // var j = 0
    // for i<len(l) {
    //   max = 0
    //   submag := make([]float64, i+1-j)
    //   subphs := make([]float64, i+1-j)
    //   var subl = l[j:i+1]
    //   // fmt.Println(len(submag))
    //   subfourier :=  fft.FFTReal(subl)
    //   for k,x :=range subfourier {
    //     submag[k],subphs[k] = cmplx.Polar(x)
    //     if submag[k] > max {
    //       max = submag[k]
    //     }
    //   }
    //   // fmt.Println(len(mag))
    //   mag = append(mag, submag...)
    //   phs = append(phs, subphs...)
    //   j=i+1
    //   i+=SAMPLE_PER_FRAME
    //   if len(l)-i>0&&len(l)-i<SAMPLE_PER_FRAME{
    //     i=len(l)-1
    //   }
    // }
    //-------------------divide into frames--------------------------

    var pi = math.Pi
    step := [5]float64{pi/48,pi/40,pi/32,pi/24,pi/16}

    var info = "Nguyen Trong Tin - Graduation Thesis - University of Science"
    var stringbit = ""
    byteArray := []byte(info)
    for _, char := range byteArray{
      n := int64(char)
      substr := strconv.FormatInt(n, 2)// 111001
      if len(substr) < 8{
        length := len(substr)
        for j:=1;j<=8-length;j++{
          substr = "0" + substr
        }
      }
      stringbit += substr
    }
    fmt.Println(stringbit)

    var k =BIT_OFFSET
    var count=0
    var pos=0
    for pos<len(stringbit){
      if math.Abs(mag[k]) < FREQ_THRES {
        k++
        continue
      }
      var stepsize = findStep(mag[k])
      if stringbit[pos] == '0'{
        phs[k] = math.Floor(phs[k]/step[stepsize] + 0.5)*step[stepsize]
      }
      if stringbit[pos] == '1' {
        phs[k] = math.Floor(phs[k]/step[stepsize])*step[stepsize] + step[stepsize]/2
      }
      // fmt.Println(phs[k])
      count++
      if count==BIT_REPEAT{
        count=0
        pos++
      }
      k++
    }
    fmt.Println(phs[103], " ",phs[104]," ", phs[105]," ",phs[107])
    cmplxArray := make([]complex128, reader.Samples)
    for i,_ := range mag {
      cmplxArray[i] = cmplx.Rect(mag[i],phs[i])
    }

    //----------------------ifft the whole file-------------------
    var wm_frame = fft.IFFT(cmplxArray)
    for i :=0 ;i<10;i++ {
      fmt.Print(real(wm_frame[i]), " ")
    }

    var newWav = make([]float64, reader.Samples)
    for i,_ := range wm_frame {
      newWav[i] = real(wm_frame[i])
    }
    //----------------------ifft the whole file-------------------

    //==================================================================

    //----------------------divide into frames----------------------
    // i = SAMPLE_PER_FRAME-1
    // j=0
    // var newWav = make([]float64, 0)
    // for i<len(l) {
    //   var subcmplx = cmplxArray[j:i+1]
    //   subIFFT :=  fft.IFFTRealOutput(subcmplx)
    //   // fmt.Println(len(newWav))
    //   newWav = append(newWav, subIFFT...)
    //   j=i+1
    //   i+=SAMPLE_PER_FRAME
    //   if len(l)-i>=0&&len(l)-i<SAMPLE_PER_FRAME{
    //     i=len(l)-1
    //   }
    // }
    //----------------------divide into frames----------------------

    wavOut, err := os.Create("test_wm.wav")
  	checkErr(err)
  	defer wavOut.Close()

  	meta := wavwriter.File{
  		Channels:        1,
  		SampleRate:      reader.Header.SampleRate,
  		SignificantBits: reader.Header.BitsPerSample,
  	}

  	writer, err := meta.NewWriter(wavOut)
  	checkErr(err)
  	defer writer.Close()

  	for n := 0; n < reader.Samples; n += 1 {
      // integer := int16(newWav[n]*math.MaxInt16)
      integer := int16(newWav[n]*(math.MaxInt16 - math.MinInt16) + math.MinInt16)
  		err = writer.WriteInt16(integer)
  		checkErr(err)
  	}
}

func findStep(mag float64) int32{
  var sMag = mag/(0.005)
  var group = math.Ceil(sMag/0.2)
  if group==0{
    group=0
  }
  if group>4{
    group=4
  }
  return int32(group)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
