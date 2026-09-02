document.querySelectorAll('[jsaction^="mouseover:"],[jsaction^="mouseout:"],button,[role="navigation"],.scene-footer-container,#omnibox-container,[aria-label="Sign in"],[aria-label="Google apps"]')
	.forEach(el => el.style.display = 'none');
document.querySelectorAll('[aria-label="Interactive map"]').forEach(el => el.parentNode.style.display='none');