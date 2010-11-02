package sqs

import "com.abneptis.oss/urltools"

func sqsEscapeTest(i byte)(out bool){
  switch i {
    case 'a','b','c','d','e','f','g','h','i','j','k','l','m',
         'A','B','C','D','E','F','G','H','I','J','K','L','M',
         'n','o','p','q','r','s','t','u','v','w','x','y','z',
         'N','O','P','Q','R','S','T','U','V','W','X','Y','Z',
         '0','1','2','3','4','5','6','7','8','9','-':
      out = false
    default:
      out = true
  }
  return
}

func sqsEscape(in string)(out string){
  return urltools.Escape(in, sqsEscapeTest, urltools.PercentUpper)
}



