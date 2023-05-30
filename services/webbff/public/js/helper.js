// const hashBrowser = val =>
//          crypto.subtle.digest('SHA-256', new TextEncoder('utf-8').encode(val))
//         .then(h => {
//            let hexes = [],
//            view = new DataView(h);
//            for (let i = 0; i < view.byteLength; i += 4)
//            hexes.push(('00000000' + view.getUint32(i).toString(16)).slice(-8));
//            return hexes.join('');
//         });

const hashBrowser = val => {
  const hashBits = sjcl.hash.sha256.hash(val);

  // Convert the hash to a hexadecimal string
  const hexHash = sjcl.codec.hex.fromBits(hashBits);

  return hexHash;
}

async function check_form(form) {


  var name = form.name.value;
  var email = form.email.value;
  var password = form.password.value;

  if (name === "" || email === "" || password === "") {
    alert("Please fill out all fields.");
    return false;
  }

  // Check if the name contains only alphabetic characters and spaces
  var name_regex = /^[a-zA-Z\s]+$/;
  if (!name_regex.test(name)) {
    alert("Name must contain only alphabetic characters and spaces.");
    return false;
  }


  // Check if the email address is valid using a regular expression
  var email_regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!email_regex.test(email)) {
    alert("Please enter a valid email address.");
    return false;
  }

  // Check if the password meets the minimum requirements
  if (!/[A-Z]/.test(password) || !/[\W_]/.test(password)) {
    alert("Password must be contain at least one capital letter and one special character");
    return false;
  }
  var pass = $('#userpass');
  let hash = await hashBrowser(pass.val());
  pass.val(hash);
  // If all checks pass, submit the form
  form.submit();
}