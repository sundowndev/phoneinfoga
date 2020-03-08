<template>
  <div>
    <h1 id="resources">Resources</h1>
    <h3 id="inderstanding-phone-numbers">Inderstanding phone numbers</h3>
    <ul>
      <li>
        <a href="http://whitepages.fr/phonesystem/"
          >whitepages.fr/phonesystem</a
        >
      </li>
      <li>
        <a
          href="https://support.twilio.com/hc/en-us/articles/223183008-Formatting-International-Phone-Numbers"
          >Formatting-International-Phone-Numbers</a
        >
      </li>
      <li>
        <a
          href="https://en.wikipedia.org/wiki/National_conventions_for_writing_telephone_numbers"
          >National_conventions_for_writing_telephone_numbers</a
        >
      </li>
    </ul>
    <h3 id="open-data">Open data</h3>
    <ul>
      <li>
        <a href="https://api.ovh.com/console/#/telephony"
          >api.ovh.com/console/#/telephony</a
        >
      </li>
      <li><a href="https://countrycode.org/">countrycode.org</a></li>
      <li>
        <a href="http://www.countryareacode.net/en/">countryareacode.net</a>
      </li>
      <li>
        <a href="http://directory.didww.com/area-prefixes"
          >directory.didww.com/area-prefixes</a
        >
      </li>
      <li><a href="http://www.numinfo.net/">numinfo.net</a></li>
      <li>
        <a href="https://gist.github.com/Goles/3196253"
          >gist.github.com/Goles/3196253</a
        >
      </li>
    </ul>
    <h2 id="footprinting">Footprinting</h2>

    <b-alert show variant="info">
      <h4 class="alert-heading">Info</h4>
      <p>
        Both free and premium resources are included. Be careful, the listing of
        a data source here does not mean it has been verified or is used in the
        tool. Data might be false. Use it as an OSINT framework.
      </p>
    </b-alert>

    <h3 id="reputation-fraud">Reputation / fraud</h3>
    <ul>
      <li>scamcallfighters.com</li>
      <li>signal-arnaques.com</li>
      <li>whosenumber.info</li>
      <li>findwhocallsme.com</li>
      <li>yellowpages.ca</li>
      <li>phonenumbers.ie</li>
      <li>who-calledme.com</li>
      <li>usphonesearch.net</li>
      <li>whocalled.us</li>
      <li>quinumero.info</li>
    </ul>
    <h3 id="disposable-numbers">Disposable numbers</h3>
    <ul>
      <li>receive-sms-online.com</li>
      <li>receive-sms-now.com</li>
      <li>hs3x.com</li>
      <li>twilio.com</li>
      <li>freesmsverification.com</li>
      <li>freeonlinephone.org</li>
      <li>sms-receive.net</li>
      <li>smsreceivefree.com</li>
      <li>receive-a-sms.com</li>
      <li>receivefreesms.com</li>
      <li>freephonenum.com</li>
      <li>receive-smss.com</li>
      <li>receivetxt.com</li>
      <li>temp-mails.com</li>
      <li>receive-sms.com</li>
      <li>receivesmsonline.net</li>
      <li>receivefreesms.com</li>
      <li>sms-receive.net</li>
      <li>pinger.com (=&gt; textnow.com)</li>
      <li>receive-a-sms.com</li>
      <li>k7.net</li>
      <li>kall8.com</li>
      <li>faxaway.com</li>
      <li>receivesmsonline.com</li>
      <li>receive-sms-online.info</li>
      <li>sellaite.com</li>
      <li>getfreesmsnumber.com</li>
      <li>smsreceiving.com</li>
      <li>smstibo.com</li>
      <li>catchsms.com</li>
      <li>freesmscode.com</li>
      <li>smsreceiveonline.com</li>
      <li>smslisten.com</li>
      <li>sms.sellaite.com</li>
      <li>smslive.co</li>
    </ul>
    <h3 id="individuals">Individuals</h3>
    <ul>
      <li>Facebook</li>
      <li>Twitter</li>
      <li>Instagram</li>
      <li>Linkedin</li>
      <li>True People</li>
      <li>Fast People</li>
      <li>Background Check</li>
      <li>Pipl</li>
      <li>Spytox</li>
      <li>Makelia</li>
      <li>IvyCall</li>
      <li>PhoneSearch</li>
      <li>411</li>
      <li>USPhone</li>
      <li>WP Plus</li>
      <li>Thats Them</li>
      <li>True Caller</li>
      <li>Sync.me</li>
      <li>WhoCallsMe</li>
      <li>ZabaSearch</li>
      <li>DexKnows</li>
      <li>WeLeakInfo</li>
      <li>OK Caller</li>
      <li>SearchBug</li>
      <li>numinfo.net</li>
    </ul>
    <h3 id="google-dork-examples">Google dork examples</h3>
    <pre
      id="__code_0"
    ><code>insubject:"+XXXXXXXXX" OR insubject:"+XXXXX" OR insubject:"XXXXX XXX XXX"
insubject:"XXXXXXXXX" OR intitle:"XXXXXXXXX"
intext:"XXXXXXXXX" AND (ext:doc OR ext:docx OR ext:odt OR ext:pdf OR ext:rtf OR ext:sxw OR ext:psw OR ext:ppt OR ext:pptx OR ext:pps OR ext:csv OR ext:txt OR ext:html)
site:"hs3x.com" "+XXXXXXXXX"
site:signal-arnaques.com intext:"XXXXXXXXX" intitle:" | Phone Fraud"
</code></pre>
  </div>
</template>

<script>
export default {
  name: "Resources"
};
</script>
