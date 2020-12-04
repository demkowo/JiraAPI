#!/bin/bash 

curl --request POST --url 'https://auth.atlassian.com/oauth/token' --header 'Content-Type: application/json' --data '{"grant_type": "authorization_code","client_id": "hPYbH5ghUKk0VTMPnkvwGDSm4B4EGApC","client_secret": "ag20pruBTFmntG8s2DnUrlI7ne6vDVvCIfgLsahiVxnoef-GgxONu7vHZhMjem_0","code": "7BgShwWxONpNYVBx","redirect_uri": "https://localhost"}'