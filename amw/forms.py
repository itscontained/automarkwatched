from django import forms


class ServerForm(forms.Form):
    servers = forms.ChoiceField(widget=forms.Select(attrs={"class": "form-control"}))


class BulkEditForm(forms.Form):
    showpkid = forms.CharField(widget=forms.TextInput(attrs={"class": "form-check-input"}))
    CHOICES = (True, False)
    silenced = forms.ChoiceField(widget=forms.RadioSelect, choices=CHOICES)
