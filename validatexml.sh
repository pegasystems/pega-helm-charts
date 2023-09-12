if xmllint --noout charts/pega/config/deploy/context.xml.tmpl >/dev/null 2>&1; then
  echo "context.xml.tmpl validation passed"
else
  echo "context.xml.tmpl formatting is incorrect"
  exit 1
fi

if xmllint --noout charts/pega/config/deploy/server.xml.tmpl >/dev/null 2>&1; then
  echo "server.xml.tmpl validation passed"
else
  echo "server.xml.tmpl formatting is incorrect"
  exit 1
fi

if xmllint --noout charts/pega/config/deploy/prconfig.xml >/dev/null 2>&1; then
  echo "deploy prconfig.xml validation passed"
else
  echo "deploy prconfig.xml formatting is incorrect"
  exit 1
fi

if xmllint --noout charts/pega/config/deploy/prlog4j2.xml >/dev/null 2>&1; then
  echo "deploy prlog4j2.xml  validation passed"
else
  echo "deploy prlog4j2.xml formatting is incorrect"
  exit 1
fi

if xmllint --noout charts/pega/config/deploy/tomcat-web.xml >/dev/null 2>&1; then
  echo "deploy tomcat-web.xml  validation passed"
else
  echo "deploy tomcat-web.xml formatting is incorrect"
  exit 1
fi

if xmllint --noout charts/pega/charts/installer/config/prlog4j2.xml >/dev/null 2>&1; then
  echo "installer prlog4j2.xml validation passed"
else
  echo "installer prlog4j2.xml formatting is incorrect"
  exit 1
fi

if xmllint --noout charts/pega/charts/installer/config/prconfig.xml.tmpl >/dev/null 2>&1; then
  echo "installer prconfig.xml.tmpl validation passed"
else
  echo "installer prconfig.xml.tmpl formatting is incorrect"
  exit 1
fi