import requests,random

uagent=[]
uagent.append("Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.0) Opera 12.14")
uagent.append("Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:26.0) Gecko/20100101 Firefox/26.0")
uagent.append("Mozilla/5.0 (X11; U; Linux x86_64; en-US; rv:1.9.1.3) Gecko/20090913 Firefox/3.5.3")
uagent.append("Mozilla/5.0 (Windows; U; Windows NT 6.1; en; rv:1.9.1.3) Gecko/20090824 Firefox/3.5.3 (.NET CLR 3.5.30729)")
uagent.append("Mozilla/5.0 (Windows NT 6.2) AppleWebKit/535.7 (KHTML, like Gecko) Comodo_Dragon/16.1.1.0 Chrome/16.0.912.63 Safari/535.7")
uagent.append("Mozilla/5.0 (Windows; U; Windows NT 5.2; en-US; rv:1.9.1.3) Gecko/20090824 Firefox/3.5.3 (.NET CLR 3.5.30729)")
uagent.append("Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US; rv:1.9.1.1) Gecko/20090718 Firefox/3.5.1")
uagent.append("Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:63.0) Gecko/20100101 Firefox/63.0")

proxy=[]
proxy.append("118.97.125.150:8080")
proxy.append("212.23.250.46:80")
proxy.append("197.149.128.190:42868")
proxy.append("87.128.41.56:80",)
proxy.append("197.149.129.252:32486")
proxy.append("159.69.211.173:3128")
proxy.append("197.149.128.190:44655")
proxy.append("196.13.208.23:8080")
proxy.append("196.13.208.22:8080")
proxy.append("82.136.122.127:80")
proxy.append("178.60.28.98:9999")
proxy.append("41.60.1.102:80")
proxy.append("212.56.139.253:80")

number = '49495363899'

dorks=[]
dorks.append('site%3Anuminfo.net+intext%3A"luciusunegbu%40gmail.com"')
dorks.append('site%3A"hs3x.com"+intext%3A"%2B61437954897"')
dorks.append('site:facebook.com intext:"%s" | "%s"' % (number,number))
dorks.append('site:twitter.com intext:"%s" | "%s"' % (number,number))
dorks.append('site:linkedin.com intext:"%s" | "%s"' % (number,number))
dorks.append('site:instagram.com intext:"%s" | "%s"' % (number,number))
dorks.append('site:whosenumber.info intext:"%s" intitle:"who called"' % number)
dorks.append('site:"hs3x.com" intext:"+%s"' % number)
dorks.append('site:"receive-sms-now.com" intext:"+%s"' % number)
dorks.append('site:"receive-sms-online.com" intext:"+%s"' % number)

for req in dorks:
    chosenProxy = random.choice(proxy)
    chosenUserAgent = random.choice(uagent)

    s = requests.Session()
    proxies = {"http": chosenProxy}
    headers = {
        'User-Agent': chosenUserAgent,
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
        'Accept-Language': 'en-us,en;q=0.5',
        'Accept-Encoding': 'gzip,deflate',
        'Accept-Charset': 'ISO-8859-1,utf-8;q=0.7,*;q=0.7',
        'Keep-Alive': '115',
        'Connection': 'keep-alive',
        'Cookie': 'Cookie: CGIC=Ij90ZXh0L2h0bWwsYXBwbGljYXRpb24veGh0bWwreG1sLGFwcGxpY2F0aW9uL3htbDtxPTAuOSwqLyo7cT0wLjg; CONSENT=YES+RE.fr+20150809-08-0; 1P_JAR=2018-11-28-14; NID=148=aSdSHJz71rufCokaUC93nH3H7lOb8E7BNezDWV-PyyiHTXqWK5Y5hsvj7IAzhZAK04-QNTXjYoLXVu_eiAJkiE46DlNn6JjjgCtY-7Fr0I4JaH-PZRb7WFgSTjiFqh0fw2cCWyN69DeP92dzMd572tQW2Z1gPwno3xuPrYC1T64wOud1DjZDhVAZkpk6UkBrU0PBcnLWL7YdL6IbEaCQlAI9BwaxoH_eywPVyS9V; SID=uAYeu3gT23GCz-ktdGInQuOSf-5SSzl3Plw11-CwsEYY0mqJLSiv7tFKeRpB_5iz8SH5lg.; HSID=AZmH_ctAfs0XbWOCJ; SSID=A0PcRJSylWIxJYTq_; APISID=HHB2bKfJ-2ZUL5-R/Ac0GK3qtM8EHkloNw; SAPISID=wQoxetHBpyo4pJKE/A2P6DUM9zGnStpIVt; SIDCC=ABtHo-EhFAa2AJrJIUgRGtRooWyVK0bAwiQ4UgDmKamfe88xOYBXM47FoL5oZaTxR3H-eOp7-rE; OTZ=4671861_52_52_123900_48_436380; OGPC=873035776-8:; OGP=-873035776:; GOOGLE_ABUSE_EXEMPTION=ID=5f05257564a6b82d:TM=1543414930:C=r:IP=165.227.163.116-:S=APGng0vsERIFITaqUYe9u-Uu9E6iJxUx-g'
    }

    r = s.get('https://www.google.com/search?q=%s&amp;gws_rd=ssl' % req, headers=headers, proxies=proxies)
    #print(r.text)

    print "%s\nRequested %s using Proxy %s\nUser-Agent: %s" % (r, req, chosenProxy, chosenUserAgent)
