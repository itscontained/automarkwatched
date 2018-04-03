from django import forms

class ServerForm(forms.Form):
    url = forms.CharField(widget=forms.TextInput(attrs={"placeholder": '10.0.10.121:32400', "class": "form-control"}))
    token = forms.CharField(widget=forms.TextInput(attrs={"placeholder": 'SDc89Xasdf9SS)8dX-D(DS', "class": "form-control"}))

class BulkEditForm(forms.Form):
    showpkid = forms.CharField(widget=forms.TextInput(attrs={"class": "form-check-input"}))
    CHOICES = ( True, False)
    silenced = forms.ChoiceField(widget=forms.RadioSelect, choices=CHOICES)
