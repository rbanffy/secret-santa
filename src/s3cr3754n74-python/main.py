#!/usr/bin/env python
import os
import random
import logging

from google.appengine.api import mail

import webapp2
import jinja2

jinja_environment = jinja2.Environment(
    loader=jinja2.FileSystemLoader(
        map(os.path.abspath,
            ('./templates',))
        )
    )

class MainHandler(webapp2.RequestHandler):
    def get(self):
        template = jinja_environment.get_template('index.html')
        self.response.write(template.render({}))

    def post(self):
        # import pdb; pdb.set_trace()
        emails = self.request.POST['emails'].split()
        random.shuffle(emails)
        sender_address = "Amigo secreto "\
            "<no-reply@banffy.com.br>"
        subject = self.request.POST['subject']
        just_test = self.request.POST.get('just_test', False) == 'true'
        email_template = jinja_environment.get_template('email.txt')
        response_template = jinja_environment.get_template('response.html')
        if just_test:
            messages = []
            for i in range(len(emails)):
                messages.append('%s will gift %s' % (emails[i-1], emails[i]))
            output = response_template.render({'emails': messages})
        else:
            for i in range(len(emails)):
                body = email_template.render({'kid': emails[i]})
                user_address = emails[i-1]
                try:
                    mail.send_mail(sender_address, user_address, subject, body)
                except mail.InvalidSenderError:
                    logging.error('Invalid sender error %s' % sender_address)
                    self.abort(500)
                    # logging.error('invalid e-mail %s' % user_address)
            output = response_template.render({'emails': emails})
        self.response.write(output)

app = webapp2.WSGIApplication([
    ('/', MainHandler)
], debug=True)
