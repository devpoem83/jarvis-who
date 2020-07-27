FROM buildpack-deps:buster-scm

ADD mapper /mapper
ADD profile /profile
ADD public /public
ADD jarvis-who /jarvis-who

CMD /jarvis-who
